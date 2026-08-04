package main

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	proxmoxURL  = "https://192.168.1.210:8006/api2/json"
	node        = "pve"
	tokenID     = "root@pam!prodqa"
	tokenSecret = "e25e44ca-34e9-4676-9b40-5ba07b706311"

	sourceVMID = 111 // your template VM
	vmCount    = 2   // number of VMs to clone
	bridge     = "vmbr0"

	maxParallel = 5 // limit concurrency
)

// Get next available VMID from Proxmox
func getNextVMID() (int, error) {
	urlAPI := proxmoxURL + "/cluster/nextid"
	req, _ := http.NewRequest("GET", urlAPI, nil)
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}
	return strconv.Atoi(result.Data)
}

// Generate random MAC address
func generateMAC() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	buf[0] = (buf[0] | 2) & 0xfe
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

// Clone VM
func cloneVM(vmid int, name string) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/clone", proxmoxURL, node, sourceVMID)
	form := url.Values{}
	form.Add("newid", strconv.Itoa(vmid))
	form.Add("name", name)
	form.Add("full", "1") // full clone

	req, _ := http.NewRequest("POST", urlAPI, bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("clone failed: %s - %s", resp.Status, string(body))
	}

	fmt.Printf("✅ Cloned VM %d\n", vmid)
	return nil
}

// Wait for VM config and lock release
func waitForVMConfig(vmid int) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/config", proxmoxURL, node, vmid)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	for i := 0; i < 60; i++ { // retry 60 times (~60 sec)
		req, _ := http.NewRequest("GET", urlAPI, nil)
		req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
		resp, err := client.Do(req)
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if resp.StatusCode == 200 && !bytes.Contains(body, []byte("lock")) {
				return nil // config exists and not locked
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("VM config %d not available after waiting", vmid)
}

// Set random MAC
func setMAC(vmid int) error {
	mac := generateMAC()
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/config", proxmoxURL, node, vmid)
	form := url.Values{}
	form.Add("net0", fmt.Sprintf("virtio=%s,bridge=%s", mac, bridge))

	req, _ := http.NewRequest("PUT", urlAPI, bytes.NewBufferString(form.Encode()))
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	fmt.Printf("🔧 MAC set for VM %d → %s\n", vmid, mac)
	return nil
}

// Start VM after lock is released
func startVM(vmid int) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/status/start", proxmoxURL, node, vmid)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	for i := 0; i < 30; i++ { // wait for lock release (~30 sec)
		req, _ := http.NewRequest("POST", urlAPI, nil)
		req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				fmt.Printf("🚀 VM %d started\n", vmid)
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to start VM %d: still locked or error", vmid)
}

func main() {
	startID, err := getNextVMID()
	if err != nil {
		fmt.Println("Error getting next VMID:", err)
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxParallel)

	for i := 0; i < vmCount; i++ {
		wg.Add(1)
		vmid := startID + i
		name := fmt.Sprintf("clone-vm-%d", vmid)

		go func(vmid int, name string) {
			defer wg.Done()
			sem <- struct{}{} // acquire semaphore

			if err := cloneVM(vmid, name); err != nil {
				fmt.Println("Clone error:", err)
				<-sem
				return
			}

			if err := waitForVMConfig(vmid); err != nil {
				fmt.Println("Config wait error:", err)
				<-sem
				return
			}

			if err := setMAC(vmid); err != nil {
				fmt.Println("Set MAC error:", err)
			}

			if err := startVM(vmid); err != nil {
				fmt.Println("Start VM error:", err)
			}

			<-sem // release semaphore
		}(vmid, name)
	}

	wg.Wait()
	fmt.Println("All VM operations completed.")
}

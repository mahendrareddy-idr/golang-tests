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
	"time"
)

const (
	proxmoxURL  = "https://192.168.1.210:8006/api2/json"
	node        = "pve" // your Proxmox node
	tokenID     = "root@pam!prodqa"
	tokenSecret = "e25e44ca-34e9-4676-9b40-5ba07b706311"

	sourceVMID = 111 // your template VM ID
	vmCount    = 1   // number of clones
	bridge     = "vmbr0"
)

// Get next available VMID
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

// Clone template VM
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

// Wait until VM config exists
func waitForVMConfig(vmid int) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/config", proxmoxURL, node, vmid)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}

	for i := 0; i < 20; i++ { // try 20 times (max ~20 sec)
		req, _ := http.NewRequest("GET", urlAPI, nil)
		req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			resp.Body.Close()
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("VM config %d not available after waiting", vmid)
}

// Set MAC address
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
	defer resp.Body.Close()

	fmt.Printf("🔧 MAC set for VM %d → %s\n", vmid, mac)
	return nil
}

// Start VM
func startVM(vmid int) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu/%d/status/start", proxmoxURL, node, vmid)
	req, _ := http.NewRequest("POST", urlAPI, nil)
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("🚀 VM %d started\n", vmid)
	return nil
}

func main() {
	startID, err := getNextVMID()
	if err != nil {
		fmt.Println("Error getting next VMID:", err)
		return
	}

	for i := 0; i < vmCount; i++ {
		vmid := startID + i
		name := fmt.Sprintf("clone-vm-%d", vmid)

		if err := cloneVM(vmid, name); err != nil {
			fmt.Println("Clone error:", err)
			continue
		}

		// Wait until VM config exists to avoid lock/config errors
		if err := waitForVMConfig(vmid); err != nil {
			fmt.Println("Config wait error:", err)
			continue
		}

		if err := setMAC(vmid); err != nil {
			fmt.Println("Set MAC error:", err)
		}

		if err := startVM(vmid); err != nil {
			fmt.Println("Start VM error:", err)
		}

		time.Sleep(1 * time.Second) // small delay to avoid lock conflicts
	}
}

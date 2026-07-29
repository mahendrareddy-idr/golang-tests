package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	proxmoxURL  = "https://192.168.1.210:8006/api2/json"
	node        = "pve"
	tokenID     = "root@pam!prodqa"
	tokenSecret = "e25e44ca-34e9-4676-9b40-5ba07b706311"
	storage     = "vmstorage6tb"
	bridge      = "vmbr0"

	vmCount = 30 // change as needed
	isoPath = "local:iso/Win_Pro_10_21H1_64Bit.ISO"
)

// Proxmox returns VMID as string
type NextIDResponse struct {
	Data string `json:"data"`
}

// 🔹 Get next available VMID
func getNextVMID() (int, error) {
	url := proxmoxURL + "/cluster/nextid"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var result NextIDResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, err
	}

	vmid, err := strconv.Atoi(result.Data)
	if err != nil {
		return 0, err
	}

	return vmid, nil
}

// 🔹 Create VM
func createVM(vmid int, name string) error {
	urlAPI := fmt.Sprintf("%s/nodes/%s/qemu", proxmoxURL, node)

	form := url.Values{}
	form.Add("vmid", fmt.Sprintf("%d", vmid))
	form.Add("name", name)
	form.Add("cores", "2")
	form.Add("sockets", "1")
	form.Add("memory", "2096")
	form.Add("scsihw", "virtio-scsi-pci")
	form.Add("sata0", fmt.Sprintf("%s:150", storage))
	form.Add("ide2", fmt.Sprintf("%s,media=cdrom", isoPath))

	// ✅ Boot from ISO first
	form.Add("boot", "order=ide2;sata0")

	form.Add("net0", fmt.Sprintf("virtio,bridge=%s", bridge))

	req, err := http.NewRequest("POST", urlAPI, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return fmt.Errorf("VM %d failed: %s - %s", vmid, resp.Status, buf.String())
	}

	fmt.Printf("✅ VM %d created\n", vmid)
	return nil
}

// 🔹 Start VM
func startVM(vmid int) {
	url := fmt.Sprintf("%s/nodes/%s/qemu/%d/status/start", proxmoxURL, node, vmid)

	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Set("Authorization", fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret))

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err == nil && resp.StatusCode == 200 {
		fmt.Printf("🚀 VM %d started\n", vmid)
	} else {
		fmt.Printf("⚠️ Failed to start VM %d\n", vmid)
	}
}

func main() {

	// 🔥 Get starting VMID once (fast + efficient)
	startID, err := getNextVMID()
	if err != nil {
		fmt.Println("Error getting VMID:", err)
		return
	}

	for i := 0; i < vmCount; i++ {
		vmid := startID + i
		name := fmt.Sprintf("Idrive-vm-%d", vmid)

		err := createVM(vmid, name)
		if err != nil {
			fmt.Println(err)
			continue
		}

		startVM(vmid)
	}
}

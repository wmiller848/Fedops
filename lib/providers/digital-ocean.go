// The MIT License (MIT)

// Copyright (c) 2014 William Miller

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package fedops_provider

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	_ "io/ioutil"
	"net/http"
	"strconv"
)

const DigitalOceanName string = "DigitalOcean"

type digitalOceanKeyRequest struct {
	Name       string `json:"name"`
	Public_key string `json:"public_key"`
}

type digitalOceanVMRequest struct {
	Name    string   `json:"name"`
	Region  string   `json:"region"`
	Size    string   `json:"size"`
	Image   string   `json:"image"`
	Keys    []string `json:"ssh_keys"`
	Backups bool     `json:"backups"`
	IPV6    bool     `json:"ipv6"`
	//UserData string `json:"user_data"`
	//PrivateNetworking bool `json:"private_networking"`
}

type DigitalOceanAuth struct {
	ApiKey string
}

type DigitalOcean struct {
	ApiKey      string
	ApiEndpoint string
	KeyURI      string
	SizeURI     string
	ImageURI    string
	VM_URI      string
}

func (digo *DigitalOcean) Name() string {
	return DigitalOceanName
}

func (digo *DigitalOcean) CreateKeypair(clusterid string, keypair Keypair) (ProviderKeypair, error) {
	client := &http.Client{}
	reqKey := digitalOceanKeyRequest{
		Name:       "FedOps-ClusterKey-" + clusterid,
		Public_key: string(keypair.PublicSSH),
	}
	reqJSON, err := json.Marshal(reqKey)
	if err != nil {
		return ProviderKeypair{}, err
	}
	req, err := http.NewRequest("POST", digo.ApiEndpoint+digo.KeyURI, bytes.NewBuffer(reqJSON))
	req.Header.Add("X-FedOps-Provider", DigitalOceanName)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+digo.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return ProviderKeypair{}, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	var data interface{}

	err = decoder.Decode(&data)
	if err != nil {
		return ProviderKeypair{}, err
	}

	jsonMap := data.(map[string]interface{})
	ssh_key := jsonMap["ssh_key"].(map[string]interface{})

	ids := make(map[string]string)
	ids[DigitalOceanName] = strconv.FormatFloat(ssh_key["id"].(float64), 'f', 0, 32)
	pkeypair := ProviderKeypair{
		ID:      ids,
		Keypair: keypair,
	}
	return pkeypair, nil
}


func (digo *DigitalOcean) ListSize() (ProviderSize, error) {
  return ProviderSize{}, nil
}

func (digo *DigitalOcean) ListSizes() ([]ProviderSize, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", digo.ApiEndpoint+digo.SizeURI, nil)
	req.Header.Add("X-FedOps-Provider", DigitalOceanName)
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+digo.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	var data interface{}

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	jsonMap := data.(map[string]interface{})
	sizes := jsonMap["sizes"].([]interface{})

	var psizes []ProviderSize
	for _, sizevalue := range sizes {
		id := sizevalue.(map[string]interface{})["slug"]
		memory := sizevalue.(map[string]interface{})["memory"]
		vcpus := sizevalue.(map[string]interface{})["vcpus"]
		disk := sizevalue.(map[string]interface{})["disk"]
		bandwidth := sizevalue.(map[string]interface{})["transfer"]
		price := sizevalue.(map[string]interface{})["price_monthly"]

		ids := make(map[string]string)
		ids[DigitalOceanName] = id.(string)
		psize := ProviderSize{
			ID:        ids,
			Memory:    memory.(float64),
			Vcpus:     vcpus.(float64),
			Disk:      disk.(float64),
			Bandwidth: bandwidth.(float64),
			Price:     price.(float64),
		}
		psizes = append(psizes, psize)
	}
	return psizes, nil
}

func (digo *DigitalOcean) GetDefaultSize() (ProviderSize, error) {
	sizes, err := digo.ListSizes()
	if err != nil {
		return sizes[0], err
	}
	for index, sizevalue := range sizes {
		if sizevalue.ID[DigitalOceanName] == "512mb" {
			return sizes[index], nil
		}
	}
	return sizes[0], nil
}

func (digo *DigitalOcean) ListImage() (ProviderImage, error) {
  return ProviderImage{}, nil
}

func (digo *DigitalOcean) ListImages() ([]ProviderImage, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", digo.ApiEndpoint+digo.ImageURI, nil)
	req.Header.Add("X-FedOps-Provider", DigitalOceanName)
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+digo.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	var data interface{}

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	jsonMap := data.(map[string]interface{})
	images := jsonMap["images"].([]interface{})

	var pimages []ProviderImage
	for _, imagevalue := range images {
		id := imagevalue.(map[string]interface{})["id"]
		distro := imagevalue.(map[string]interface{})["distribution"]
		slug := imagevalue.(map[string]interface{})["slug"]

		ids := make(map[string]string)
		ids[DigitalOceanName] = strconv.FormatFloat(id.(float64), 'f', 0, 32)
		pimage := ProviderImage{
			ID:           ids,
			Distribution: distro.(string),
			Version:      slug.(string),
		}
		pimages = append(pimages, pimage)
	}
	return pimages, nil
}

func (digo *DigitalOcean) GetDefaultImage() (ProviderImage, error) {
	images, err := digo.ListImages()
	if err != nil {
		return images[0], err
	}
	for index, imagevalue := range images {
		if imagevalue.Version == "fedora-20-x64" {
			return images[index], nil
		}
	}
	return images[0], nil
}

func (digo *DigitalOcean) ListVM() (ProviderVM, error) {
  return ProviderVM{}, nil
}

func (digo *DigitalOcean) ListVMs() ([]ProviderVM, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", digo.ApiEndpoint+digo.VM_URI, nil)
	req.Header.Add("X-FedOps-Provider", DigitalOceanName)
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+digo.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	var data interface{}

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	jsonMap := data.(map[string]interface{})
	droplets := jsonMap["droplets"].([]interface{})

	var pvms []ProviderVM
	for _, vmvalue := range droplets {
		id := vmvalue.(map[string]interface{})["id"]
		networks := vmvalue.(map[string]interface{})["networks"].(map[string]interface{})
		ipv4 := networks["v4"].([]interface{})
		v4 := ipv4[0].(map[string]interface{})
		//ipv6 := networks["v6"].([]interface{})
		//v6 := ipv6[0].(map[string]interface{})

		ids := make(map[string]string)
		ids[DigitalOceanName] = strconv.FormatFloat(id.(float64), 'f', 0, 32)
		pvm := ProviderVM{
			ID:   ids,
			IPV4: v4["ip_address"].(string),
			//IPV6: v6["ip_address"].(string),
			Provider: DigitalOceanName,
		}
		pvms = append(pvms, pvm)
	}
	return pvms, nil
}

func (digo *DigitalOcean) CreateVM(vmid string, size ProviderSize, image ProviderImage, keypairs []ProviderKeypair) (ProviderVM, error) {

	keyids := []string{}
	for _, keypair := range keypairs {
		keyids = append(keyids, keypair.ID[DigitalOceanName])
	}

	client := &http.Client{}
	reqVM := digitalOceanVMRequest{
		Name:    "FedOpsWarehouse-" + vmid,
		Region:  "nyc2",
		Size:    size.ID[DigitalOceanName],
		Image:   image.ID[DigitalOceanName],
		Keys:    keyids,
		Backups: false,
		IPV6:    true,
	}
	reqJSON, err := json.Marshal(reqVM)
	if err != nil {
		return ProviderVM{}, err
	}

	req, err := http.NewRequest("POST", digo.ApiEndpoint+digo.VM_URI, bytes.NewBuffer(reqJSON))
	req.Header.Add("X-FedOps-Provider", DigitalOceanName)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+digo.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
		return ProviderVM{}, err
	}
	defer resp.Body.Close()

	//fmt.Println("Response Status:", resp.Status)
	//fmt.Println("Response Headers:", resp.Header)
	decoder := json.NewDecoder(resp.Body)
	var data interface{}

	err = decoder.Decode(&data)
	if err != nil {
		return ProviderVM{}, err
	}

	jsonMap := data.(map[string]interface{})
	droplet := jsonMap["droplet"].(map[string]interface{})

	ids := make(map[string]string)
	ids[DigitalOceanName] = strconv.FormatFloat(droplet["id"].(float64), 'f', 0, 32)
	pvm := ProviderVM{
		ID:       ids,
		Provider: DigitalOceanName,
	}
	return pvm, nil
}

func (digo *DigitalOcean) SnapShotVM(ProviderVM) (ProviderImage, error) {
	return ProviderImage{}, nil
}

func DigitalOceanProvider(auth DigitalOceanAuth) DigitalOcean {
	return DigitalOcean{
		ApiKey:      auth.ApiKey,
		ApiEndpoint: "https://api.digitalocean.com",
		KeyURI:      "/v2/account/keys",
		SizeURI:     "/v2/sizes",
		ImageURI:    "/v2/images?type=distribution",
		VM_URI:      "/v2/droplets",
	}
}

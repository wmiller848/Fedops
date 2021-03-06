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

package fedops

import (
	// Standard
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	// 3rd Party
	// FedOps
	"github.com/wmiller848/Fedops/lib/encryption"
	"github.com/wmiller848/Fedops/lib/engine/container"
	"github.com/wmiller848/Fedops/lib/engine/truck"
	"github.com/wmiller848/Fedops/lib/engine/warehouse"
	"github.com/wmiller848/Fedops/lib/providers"
)

const (
	DigitalOcean   uint = 0
	AWS            uint = 1
	GoogleCloud    uint = 2
	MicrosoftAzure uint = 3
	OpenStack      uint = 4

	SaltSize int = 512

	ClusterIDSize   int = 32
	WarehouseIDSize int = 16
	TruckIDSize     int = 16
	ContainerIDSize int = 16

	AuthorizationCost int = 10

	FedopsRemoteKeySize int = 32

	FedopsPoolTime     time.Duration = 5
	FedopsBootWaitTime time.Duration = 30

	FedopsError   uint = 0
	FedopsOk      uint = 1
	FedopsUnknown uint = 2

	FedopsTypeTruck     uint = 0
	FedopsTypeWarehouse uint = 1

	FedopsRepo     string = "github.com/wmiller848/Fedops"
	ConfigFileName string = "Fedops"
)

type FedopsAction struct {
	Status uint
}

type Tokens struct {
	AccessToken   string
	SecurityToken string
}

// type Services struct {
// 	ID   string
// 	Name string
// 	Repo string
// }

//
//
type DispatcherError struct {
	msg string
}

func (err *DispatcherError) Error() string {
	return err.msg
}

func (err *DispatcherError) setMsg(msg string) {
	err.msg = msg
}

//
// This config is stored encrypted on disk
type DispatcherConfig struct {
	ClusterID  string
	Created    string
	Modified   string
	Certs      []fedops_encryption.Cert
	SSHKeys    []fedops_provider.ProviderKeypair
	Tokens     map[string]Tokens
	Warehouses map[string]*fedops_warehouse.Warehouse
	Trucks     map[string]*fedops_truck.Truck
	Containers map[string]*fedops_container.Container
}

type Dispatcher struct {
	Cipherkey      []byte
	Salt           []byte
	Version        string
	PowerDirectory string
	Timeout        time.Duration
	Config         DispatcherConfig
}

func CreateDispatcher(key []byte, pwd string, session bool) (*Dispatcher, error) {

	var salt, cipherkey []byte
	if session == true {
		cipherkey = key
	} else {
		var err error
		salt, err = GetSalt(pwd)
		if err != nil {
			salt, err = fedops_encryption.GenerateRandomBytes(SaltSize)
			if err != nil {
				return nil, err
			}
		}
		cipherkey = make([]byte, len(salt)+len(key))
		cipherkey = append(cipherkey, salt...)
		cipherkey = append(cipherkey, key...)
		cipherkey = fedops_encryption.Hashkey(cipherkey)
	}

	config, err := loadConfig(cipherkey, pwd)
	if err != nil {
		return nil, err
	}

	d := &Dispatcher{
		Cipherkey:      fedops_encryption.Encode(cipherkey),
		Salt:           salt,
		Config:         config,
		Version:        "0.0.1",
		PowerDirectory: pwd,
		Timeout:        60,
	}
	return d, nil
}

func HasConfigFile(pwd string) bool {
	_, err := os.Stat(pwd + "/" + ConfigFileName)
	if err != nil {
		return false
	}
	return true
}

func GetConfigFile(pwd string) ([]byte, error) {
	return ioutil.ReadFile(pwd + "/" + ConfigFileName)
}

func GetSalt(pwd string) ([]byte, error) {
	return ioutil.ReadFile(pwd + "/.fedops-salt")
}

func loadConfig(cipherkey []byte, pwd string) (DispatcherConfig, error) {
	var config DispatcherConfig
	cdata, err := GetConfigFile(pwd)
	if err != nil {
		//  We couldn't find the config file :(
		//fmt.Println(err.Error())
		cid, err := fedops_encryption.GenerateRandomHex(ClusterIDSize)
		if err != nil {
			return config, err
		}
		config = DispatcherConfig{
			ClusterID: cid,
		}
		return config, nil
	}

	// We found the config, now unencrypt it, base64 decode it, and then marshal from json
	decrypted, err := fedops_encryption.Decrypt(cipherkey, cdata)
	if err != nil {
		return config, err
	}
	decoder := json.NewDecoder(bytes.NewBuffer(decrypted))
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

//
//
func (d *Dispatcher) error() DispatcherError {
	return DispatcherError{}
}

func (d *Dispatcher) Info() {
	fmt.Println("[WARNING] Fedops encrypts all information you provide to it...")
	fmt.Println("[WARNING] Fedops data is UNRECOVERABLE without knowning the encryption key")
}

func (d *Dispatcher) writeKeypair(sshKey fedops_encryption.Keypair, provider fedops_provider.Provider) {
	//fmt.Println(d.PowerDirectory)
	ioutil.WriteFile(d.PowerDirectory+"/"+provider.Name()+"_id_rsa.pub", sshKey.PublicPem, os.ModePerm)
	ioutil.WriteFile(d.PowerDirectory+"/"+provider.Name()+"_id_rsa", sshKey.PrivatePem, os.ModePerm)
}

func (d *Dispatcher) Unload() bool {

	now := time.Now()
	d.Config.Modified = now.UTC().String()

	pwd := d.PowerDirectory
	disjson, err := json.Marshal(d.Config)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	cipherkey, err := fedops_encryption.Decode(d.Cipherkey)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	encrypted, err := fedops_encryption.Encrypt(cipherkey, disjson)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = ioutil.WriteFile(pwd+"/"+ConfigFileName, encrypted, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	if len(d.Salt) > 0 {
		err = ioutil.WriteFile(pwd+"/.fedops-salt", d.Salt, 0666)
		if err != nil {
			fmt.Println(err.Error())
			return false
		}
	}
	return true
}

//
func (d *Dispatcher) InitCloudProvider(promise chan FedopsAction, provider string, tokens Tokens) {
	// digital-ocean, aws, google-cloud, microsoft-azure
	switch provider {
	case "digital ocean":
		auth := fedops_provider.DigitalOceanAuth{
			ApiKey: tokens.AccessToken,
		}
		digo := fedops_provider.DigitalOceanProvider(auth)
		d.Config.Tokens = make(map[string]Tokens)
		d.Config.Warehouses = make(map[string]*fedops_warehouse.Warehouse)
		d.Config.Trucks = make(map[string]*fedops_truck.Truck)
		d.Config.Containers = make(map[string]*fedops_container.Container)
		d.Config.Tokens[fedops_provider.DigitalOceanName] = tokens
		promise <- FedopsAction{
			Status: d._initProvider(&digo),
		}
	case "aws":
		fmt.Println("No API Driver... consider forking and submiting a PR")
		promise <- FedopsAction{
			Status: FedopsError,
		}
	case "google cloud":
		fmt.Println("No API Driver... consider forking and submiting a PR")
		promise <- FedopsAction{
			Status: FedopsError,
		}
	case "microsoft azure":
		fmt.Println("No API Driver... consider forking and submiting a PR")
		promise <- FedopsAction{
			Status: FedopsError,
		}
	case "openstack":
		fmt.Println("No API Driver... consider forking and submiting a PR")
		promise <- FedopsAction{
			Status: FedopsError,
		}
	default:
		fmt.Println("Unknown provider " + provider)
		promise <- FedopsAction{
			Status: FedopsError,
		}
	}
	//
	go func() {
		time.Sleep(d.Timeout * time.Second)
		// Signal to finish
		promise <- FedopsAction{
			Status: FedopsError,
		}
	}()
}

func (d *Dispatcher) _initProvider(provider fedops_provider.Provider) uint {

	// certConfig := fedops_encryption.Cert_Config{}
	// cert := fedops_encryption.GenerateCert(certConfig)

	// d.Config.Certs = append(d.Config.Certs, cert)

	keypairConfig := fedops_encryption.Keypair_Config{}
	sshKey := fedops_encryption.GenerateKeypair(keypairConfig)

	keypair, err := provider.CreateKeypair(d.Config.ClusterID, sshKey)
	if err != nil {
		fmt.Println(err.Error())
		return FedopsError
	}
	d.Config.SSHKeys = append(d.Config.SSHKeys, keypair)
	now := time.Now()
	d.Config.Created = now.UTC().String()
	d.Config.Modified = now.UTC().String()

	persisted := d.Unload()
	if persisted != true {
		return FedopsError
	}
	return FedopsOk
}

func (d *Dispatcher) Refresh(promise chan FedopsAction) {

	// Cycle through all the provider tokens
	for name, token := range d.Config.Tokens {
		switch name {
		case fedops_provider.DigitalOceanName:
			auth := fedops_provider.DigitalOceanAuth{
				ApiKey: token.AccessToken,
			}
			provider := fedops_provider.DigitalOceanProvider(auth)
			status := d._refresh(&provider)
			if status == FedopsError {
				promise <- FedopsAction{
					Status: FedopsError,
				}
				return
			}
		}
	}

	persisted := d.Unload()
	if persisted != true {
		promise <- FedopsAction{
			Status: FedopsError,
		}
	}
	promise <- FedopsAction{
		Status: FedopsOk,
	}
}

func (d *Dispatcher) _refresh(provider fedops_provider.Provider) uint {
	vms, err := provider.ListVMs()
	if err != nil {
		fmt.Println(err.Error())
		return FedopsError
	}

	warehouses := d.Config.Warehouses
	for wIndex, _ := range warehouses {
		for vIndex, _ := range vms {
			if vms[vIndex].ID[provider.Name()] == warehouses[wIndex].ID[provider.Name()] {
				warehouses[wIndex].IPV4 = vms[vIndex].IPV4
				warehouses[wIndex].Status = vms[vIndex].Status
			}
		}
	}

	trucks := d.Config.Trucks
	for tIndex, _ := range trucks {
		for vIndex, _ := range vms {
			if vms[vIndex].ID[provider.Name()] == trucks[tIndex].ID[provider.Name()] {
				trucks[tIndex].IPV4 = vms[vIndex].IPV4
				trucks[tIndex].Status = vms[vIndex].Status
			}
		}
	}

	return FedopsOk
}

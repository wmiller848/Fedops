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

package fedops_runtime

import (
	// Standard
	"bytes"
	"crypto/tls"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"time"
	// 3rd Party
	"golang.org/x/crypto/bcrypt"
	// FedOps
	"github.com/wmiller848/Fedops/lib/encryption"
	"github.com/wmiller848/Fedops/lib/engine/container"
	"github.com/wmiller848/Fedops/lib/engine/network"
)

const (
	KeyFileName    string = ".fedops-key"
	ConfigFileName string = "Fedops-Runtime"
)

type FedopsEventHandle func(event *FedopsEvent)

type FedopsEvent struct {
	ID         string
	Handle     FedopsEventHandle
	Persistant bool
	Time       time.Time
}

//
//
type RuntimeError struct {
	msg string
}

func (err RuntimeError) Error() string {
	return err.msg
}

//
// This config is stored encrypted on disk
type ClusterConfig struct {
	ClusterID  string
	MemberID   string
	Created    string
	Modified   string
	Cert       fedops_encryption.Cert
	Containers map[string]fedops_container.Container
}

type Runtime struct {
	Cipherkey      []byte
	Version        string
	PowerDirectory string
	Config         ClusterConfig
	Routes         []fedops_network.FedopsRoute
	Events         []FedopsEvent
}

func (r *Runtime) Configure(pwd string) error {
	if !r.HasKeyFile(pwd) {
		return RuntimeError{msg: "No key file located in " + pwd}
	}

	if !r.HasConfigFile(pwd) {
		return RuntimeError{msg: "No config file located in " + pwd}
	}

	cipherkey, err := r.GetKeyFile(pwd)
	if err != nil {
		//  We couldn't find the key file :(
		return err
	}

	config, err := r.Load(cipherkey, pwd)
	if err != nil {
		return err
	}

	r.Cipherkey = fedops_encryption.Encode(cipherkey)
	r.Config = config
	r.Version = "0.0.1"
	r.PowerDirectory = pwd

	if r.Config.Containers == nil {
		r.Config.Containers = make(map[string]fedops_container.Container)
	}

	return nil
}

func (r *Runtime) HasKeyFile(pwd string) bool {
	_, err := os.Stat(pwd + "/" + KeyFileName)
	if err != nil {
		return false
	}
	return true
}

func (r *Runtime) GetKeyFile(pwd string) ([]byte, error) {
	return ioutil.ReadFile(pwd + "/" + KeyFileName)
}

func (r *Runtime) HasConfigFile(pwd string) bool {
	_, err := os.Stat(pwd + "/" + ConfigFileName)
	if err != nil {
		return false
	}
	return true
}

func (r *Runtime) GetConfigFile(pwd string) ([]byte, error) {
	return ioutil.ReadFile(pwd + "/" + ConfigFileName)
}

//
//
func (r *Runtime) Error() *RuntimeError {
	return &RuntimeError{}
}

func (r *Runtime) Info() {
	fmt.Println("[WARNING] Fedops encrypts all information you provide to it...")
	fmt.Println("[WARNING] Fedops data is UNRECOVERABLE without knowning the encryption key")
}

func (r *Runtime) Load(cipherkey []byte, pwd string) (ClusterConfig, error) {
	var config ClusterConfig
	cdata, err := r.GetConfigFile(pwd)
	if err != nil {
		//  We couldn't find the config file :(
		return config, err
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

func (r *Runtime) Unload() bool {
	now := time.Now()
	r.Config.Modified = now.UTC().String()

	pwd := r.PowerDirectory
	disjson, err := json.Marshal(r.Config)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	cipherkey, err := fedops_encryption.Decode(r.Cipherkey)
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
	return true
}

func (r *Runtime) UnloadToMemory() []byte {
	now := time.Now()
	r.Config.Modified = now.UTC().String()

	disjson, err := json.Marshal(r.Config)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	cipherkey, err := fedops_encryption.Decode(r.Cipherkey)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	encrypted, err := fedops_encryption.Encrypt(cipherkey, disjson)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return encrypted
}

func (r *Runtime) AddRoute(method uint, route string, handle fedops_network.HandleRoute) error {
	rgx, err := regexp.Compile(route)
	if err != nil {
		return err
	}
	fedRoute := fedops_network.FedopsRoute{
		Method: method,
		Route:  rgx,
		Handle: handle,
	}
	r.Routes = append(r.Routes, fedRoute)
	return nil
}

// Handles incoming requests.
func (r *Runtime) HandleConnection(conn net.Conn) {
	// Make a buffer to hold incoming data.
	// buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	// reqLen, err := conn.Read(buf)
	// if err != nil {
	//   fmt.Println("Error reading:", err.Error())
	//   return
	// }
	// if reqLen > 0 {
	//   fmt.Println(buf)
	// }
	// Send a response back to person contacting us.
	// conn.Write([]byte("Message received"))
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	var req fedops_network.FedopsRequest
	err := dec.Decode(&req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	res := fedops_network.FedopsResponse{
		Success: true,
	}

	err = bcrypt.CompareHashAndPassword(req.Authorization, []byte(r.Config.ClusterID))
	if err != nil {
		fmt.Println("Authorization not accepted", err.Error())
		return
	} else {
		fmt.Println("Authorization accepted")
		fmt.Println("Method", req.Method)
		fmt.Println("Route", string(req.Route))
		fmt.Println("Data", string(req.Data))

		for i := range r.Routes {
			if r.Routes[i].Method == req.Method && r.Routes[i].Route.Match(req.Route) {
				err = r.Routes[i].Handle(&req, &res)
				if err != nil {
					res.Success = false
					res.Error = []byte(err.Error())
					fmt.Println(err.Error())
				}
				break
			}
		}
	}

	enc := gob.NewEncoder(conn)
	err = enc.Encode(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	persisted := r.Unload()
	if persisted != true {
		fmt.Println("Error saving to disk")
	}
	// conn.Write([]byte("ok"))
}

func (r *Runtime) Listen(status chan error) {
	fed_cert := r.Config.Cert
	// cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	cert, err := tls.X509KeyPair(fed_cert.CertificatePem, fed_cert.PrivatePem)
	if err != nil {
		fmt.Println(err.Error())
		status <- err
	}

	config := tls.Config{
		Certificates:             []tls.Certificate{cert},
		PreferServerCipherSuites: true,
		SessionTicketsDisabled:   true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		CurvePreferences: []tls.CurveID{tls.CurveP521},
		MinVersion:       tls.VersionTLS12,
		MaxVersion:       tls.VersionTLS12,
	}
	listener, err := tls.Listen("tcp", ":13371", &config)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Println(conn.RemoteAddr(), "Connected")
		go r.HandleConnection(conn)
	}
}

func (r *Runtime) StartEventEngine(status chan error) {
	for {
		l := len(r.Events)
		if l > 0 {
			fmt.Println("Processing Event")
			event := r.Events[l-1 : l][0]
			ftime := event.Time.Add(2 * time.Second)
			n := time.Now()
			if n.After(ftime) {
				fmt.Println("Calling Handle for Event")
				event.Time = n
				go event.Handle(&event)
			}
		}
	}
}

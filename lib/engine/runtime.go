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
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "time"
  // 3rd Party
  // FedOps
  "github.com/Fedops/lib/encryption"
  "github.com/Fedops/lib/engine/container"
)

const (
  KeyFileName string =  ".fedops-key"
  ConfigFileName string = "Fedops-Runtime"
)

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
  Certs      []fedops_encryption.Cert
  Containers []fedops_container.Container
}

type Runtime struct {
  Cipherkey      []byte
  Version        string
  PowerDirectory string
  Config         ClusterConfig
}

func (r *Runtime) Configure(pwd string) error {
  if !r.HasKeyFile(pwd) {
    return RuntimeError{msg:"No key file located in " + pwd}
  }

  if !r.HasConfigFile(pwd) {
    return RuntimeError{msg:"No config file located in " + pwd}
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

  err = ioutil.WriteFile(pwd+"/" + ConfigFileName, encrypted, 0666)
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

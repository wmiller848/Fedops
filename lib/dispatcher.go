package fedops

import (
  // Standard
  "os"
  "io/ioutil"
  "time"
  "fmt"
  // 3rd Party
  // FedOps
  "github.com/FedOps/lib/providers"
)

type DispatcherConfig struct {
  ClusterID string
  Created string
  Modified string
  Keys ProviderKeys
  VMs []VM
}

type Dispatcher struct {
  Verison string
  PowerDirectory string
  Timeout time.Duration
  Config DispatcherConfig
  Error uint
  Ok uint
  Unknown uint
}

func CreateDispatcher(key, pwd string) *Dispatcher {

  config, err := load(pwd)
  if err != nil {
    fmt.Println(err.Error())
    return nil
  }

  d := &Dispatcher {
    Config: config,
    Verison: "0.0.1",
    PowerDirectory: pwd,
    Timeout: 60,
    Error: 0,
    Ok: 1,
    Unknown: 2,
  }
  fmt.Printf("%+v \r\n", d)
  return d
}

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
//
func (d *Dispatcher) error() (DispatcherError) {
  return DispatcherError{}
}

func (d *Dispatcher) decrypt(bytz []byte) []byte {
  return []byte("jibber_jabber")
}

func (d *Dispatcher) encrypt(bytz []byte) []byte {
  return []byte("jibber_jabber")
}

func (d *Dispatcher) decode(bytz []byte) []byte {
  return []byte("jibber_jabber")
}

func (d *Dispatcher) encode(bytz []byte) []byte {
  return []byte("jibber_jabber")
}

func (d *Dispatcher) Unload() bool {
  return true
}

func load(pwd string) (DispatcherConfig, error) {
  fdata, err := ioutil.ReadFile(pwd + "/.fedops")
  var config DispatcherConfig
  if err != nil {
    //  We couldn't find the encrypted config file :(
    //fmt.Println(err.Error())
    cid, err := GenerateRandomString(128)
    if err != nil {
      fmt.Println(err.Error())
      return config, err
    }
    config = DispatcherConfig{
      ClusterID: cid,
    }
  } else {
    // We found the config, now unecrypt it, base64 decode it, and then marshal from json
    fmt.Println(fdata)
  }
  return config, nil
}

func (d *Dispatcher) Info() () {
  fmt.Println("[WARNING] Fedops encrypts all information you provide to it...")
  fmt.Println("[WARNING] Fedops data is UNRECOVERABLE without knowning the encryption key")
}

//
//
func (d *Dispatcher) InitCloudProvider(promise chan uint, provider string, providerTokens ProviderTokens) {
  //
  // digital-ocean, aws, google-cloud, microsoft-azure
  go func() {
    switch provider {
      case "digital ocean":
        promise <- d.initDigitalOcean(providerTokens)
      case "aws":
        fmt.Println("No API Driver :(")
        promise <- d.Error
      case "google cloud":
        fmt.Println("No API Driver :(")
        promise <- d.Error
      case "microsoft azure":
        fmt.Println("No API Driver :(")
        promise <- d.Error
      default:
        fmt.Println("Unknown provider " + provider)
        promise <- d.Error
    }
    
    go func() {
      time.Sleep(d.Timeout * time.Second)
      // Signal to finish
      promise <- d.Error
    }()
  }()
}

func (d *Dispatcher) createKeypair(sshKeyConfig fedops_provider.SSH_Config) (fedops_provider.Keypair) {
  sshkey := fedops_provider.Keypair{Keysize: sshKeyConfig.Keysize}
  sshkey.Generate()
  return sshkey
}

// REMOVE
func (d *Dispatcher) writeKeypair(sshKey fedops_provider.Keypair, provider string) {
  //fmt.Println(d.PowerDirectory)
  ioutil.WriteFile(d.PowerDirectory + "/" + provider + "_id_rsa.pub", sshKey.PublicPem, os.ModePerm)
  ioutil.WriteFile(d.PowerDirectory + "/" + provider + "_id_rsa", sshKey.PrivatePem, os.ModePerm)
}

func (d *Dispatcher) initDigitalOcean(providerTokens ProviderTokens) (uint) {
  //d.Info()
  sshKeyConfig := fedops_provider.SSH_Config { Keysize: 4096 }
  //fmt.Println("Generating new ssh keys for cluster at " + strconv.FormatInt(int64(sshKeyConfig.Keysize), 10) + " bytes")
  sshKey := d.createKeypair(sshKeyConfig)
  auth := fedops_provider.DigitalOceanAuth {
    ApiKey: providerTokens.AccessToken,
  }
  //
  //
  digo := fedops_provider.DigitalOceanProvider(auth)
  keyid, err := digo.CreateKeypair(sshKey)
  if err != nil {
    fmt.Println(err.Error())
    return d.Error
  }
  keyMap := make(map[string]fedops_provider.Keypair)
  keyMap[keyid] = sshKey
  d.Config.Keys = ProviderKeys{
    DigitalOcean: keyMap,
  }
  
  //d.writeKeypair(sshKey, "digital_ocean")
  
  fmt.Printf("%+v \r\n", d)
  return d.Ok
}

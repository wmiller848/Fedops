package fedops

import (
  // Standard
  "os"
  "io/ioutil"
  "bytes"
  "time"
  "fmt"
  "encoding/json"
  // 3rd Party
  // FedOps
  "github.com/FedOps/lib/providers"
)

const (
  DigitalOcean uint = 0
  AWS uint = 1
  GoogleCloud uint = 2
  MicrosoftAzure uint = 3
  OpenStack uint = 4
)

type ProviderTokens struct {
  AccessToken string
  SecurityToken string
}

type VM struct {
  Provider string
  Role uint
  IP string
  Aliases []string
}

//
//
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
type DispatcherConfig struct {
  ClusterID string
  //Users []User
  Created string
  Modified string
  Keys map[string]fedops_provider.Keypair
  Tokens map[string][]ProviderTokens
  VMs []VM
}

type Dispatcher struct {
  Version string
  PowerDirectory string
  Timeout time.Duration
  Config DispatcherConfig
  Error uint
  Ok uint
  Unknown uint
}

func CreateDispatcher(key, pwd string) (*Dispatcher, bool) {

  config, loaded, err := load(pwd)
  if err != nil {
    fmt.Println(err.Error())
    return nil, loaded
  }

  d := &Dispatcher {
    Config: config,
    Version: "0.0.1",
    PowerDirectory: pwd,
    Timeout: 60,
    Error: 0,
    Ok: 1,
    Unknown: 2,
  }
  //fmt.Printf("%+v \r\n", d)
  return d, loaded
}

func GetConfigFile(pwd string) ([]byte, error) {
  return ioutil.ReadFile(pwd + "/.fedops")
  
}

func load(pwd string) (DispatcherConfig, bool, error) {
  fdata, err := GetConfigFile(pwd)
  var config DispatcherConfig
  if err != nil {
    //  We couldn't find the encrypted config file :(
    //fmt.Println(err.Error())
    cid, err := GenerateRandomString(128)
    if err != nil {
      fmt.Println(err.Error())
      return config, false, err
    }
    config = DispatcherConfig {
      ClusterID: cid,
    }
    return config, false, nil
  }

  // We found the config, now unecrypt it, base64 decode it, and then marshal from json
  decrypted := decrypt(fdata)
  debased := decode(decrypted)
  decoder := json.NewDecoder(bytes.NewBuffer(debased))
  err = decoder.Decode(&config)
  if err != nil {
    fmt.Println(err.Error())
    return config, true, err;
  }

  return config, true, nil
}

//
//
func (d *Dispatcher) error() (DispatcherError) {
  return DispatcherError{}
}

func (d *Dispatcher) Info() () {
  fmt.Println("[WARNING] Fedops encrypts all information you provide to it...")
  fmt.Println("[WARNING] Fedops data is UNRECOVERABLE without knowning the encryption key")
}

func (d *Dispatcher) writeKeypair(sshKey fedops_provider.Keypair, provider fedops_provider.Provider) {
  //fmt.Println(d.PowerDirectory)
  ioutil.WriteFile(d.PowerDirectory + "/" + provider.Name() + "_id_rsa.pub", sshKey.PublicPem, os.ModePerm)
  ioutil.WriteFile(d.PowerDirectory + "/" + provider.Name() + "_id_rsa", sshKey.PrivatePem, os.ModePerm)
}

func (d *Dispatcher) Unload() bool {
  pwd := d.PowerDirectory
  fmt.Printf("%+v \r\n", d)
  disjson, err := json.Marshal(d.Config)
  if err != nil {
    fmt.Println(err.Error())
    return false
  }

  err = ioutil.WriteFile(pwd + "/.fedops", disjson, os.ModePerm)
  if err != nil {
    fmt.Println(err.Error())
    return false
  }
  return true
}

//
//
func (d *Dispatcher) InitCloudProvider(promise chan uint, provider string, providerTokens ProviderTokens) {
  //
  // digital-ocean, aws, google-cloud, microsoft-azure
  go func() {
    switch provider {
      case "digital ocean":
        auth := fedops_provider.DigitalOceanAuth {
          ApiKey: providerTokens.AccessToken,
        }
        digo := fedops_provider.DigitalOceanProvider(auth)
        d.Config.Tokens = make(map[string][]ProviderTokens)
        d.Config.Tokens[digo.Name()] = make([]ProviderTokens, 0)
        d.Config.Tokens[digo.Name()] = append(d.Config.Tokens[digo.Name()], providerTokens)
        promise <- d._initProvider(&digo)
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

func (d *Dispatcher) _initProvider(provider fedops_provider.Provider) (uint) {

  sshKeyConfig := fedops_provider.SSH_Config { Keysize: 4096 }
  sshKey := fedops_provider.GenerateKeypair(sshKeyConfig)

  keyid, err := provider.CreateKeypair(sshKey)
  if err != nil {
    fmt.Println(err.Error())
    return d.Error
  }

  keyMap := make(map[string]fedops_provider.Keypair)
  keyMap[provider.Name() + "-" + keyid] = sshKey
  d.Config.Keys = keyMap
  now := time.Now()
  d.Config.Created = now.UTC().String()
  d.Config.Modified = now.UTC().String()

  persisted := d.Unload()
  if persisted != true {
    return d.Error
  }
  return d.Ok
}

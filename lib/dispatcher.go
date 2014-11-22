package fedops

import (
  // Standard
  _"bufio"
  _"os"
  "ioutil"
  "fmt"
  "time"
  _"strconv"
  // 3rd Party
  // FedOps
  "github.com/FedOps/lib/providers"
)

type Dispatcher struct {
  Verison string
  PowerDirectory string
  Timeout time.Duration
  Error uint
  Ok uint
  Unknown uint
}

func CreateDispatcher(pwd string) *Dispatcher {
  d := &Dispatcher {
    Verison: "1.0.0",
    PowerDirectory: pwd,
    Timeout: 60,
    Error: 0,
    Ok: 1,
    Unknown: 2,
  }
  return d
}

func (d *Dispatcher) Info() () {
  fmt.Println("[WARNING] Fedops encrypts all information you provide to it...")
  fmt.Println("[WARNING] Fedops data is UNRECOVERABLE without knowning the encryption key")
}

func (d *Dispatcher) InitCloudProvider(promise chan uint, provider string) {
  //
  // digital-ocean, aws, google-cloud, microsoft-azure
  switch provider {
    case "digital ocean":
      go func() {
        time.Sleep(1 * time.Millisecond)
        promise <- d.initDigitalOcean()
      }()   
    case "aws":
      fmt.Println("No API Driver :(")
      go func() {
        time.Sleep(1 * time.Millisecond)
        promise <- d.Error
      }()
    case "google cloud":
      fmt.Println("No API Driver :(")
      go func() {
        time.Sleep(1 * time.Millisecond)
        promise <- d.Error
      }()
    case "microsoft azure":
      fmt.Println("No API Driver :(")
      go func() {
        time.Sleep(1 * time.Millisecond)
        promise <- d.Error
      }()
    default:
      fmt.Println("Unknown provider " + provider)
      go func() {
        time.Sleep(1 * time.Millisecond)
        promise <- d.Error
      }()
  }
  
  go func() {
    time.Sleep(d.Timeout * time.Second)
    // Signal to finish
    promise <- d.Error
  }()
}

func (d *Dispatcher) createKeypair(sshKeyConfig fedops_provider.SSH_Config) (fedops_provider.Keypair) {
  sshkey := fedops_provider.Keypair{Keysize: sshKeyConfig.Keysize}
  sshkey.Generate()
  return sshkey
}

func (d *Dispatcher) writeKeypair(fedops_provider.Keypair) {
  ioutil.WriteFile(d.PowerDirectory + ".fedops", sshKey.Keypair)
}

func (d *Dispatcher) initDigitalOcean() (uint) {
  //d.Info()
  sshKeyConfig := fedops_provider.SSH_Config { Keysize: 4096 }
  //fmt.Println("Generating new ssh keys for cluster at " + strconv.FormatInt(int64(sshKeyConfig.Keysize), 10) + " bytes")
  sshKey := d.createKeypair(sshKeyConfig)
  auth := fedops_provider.DigitalOceanAuth{}
  digo := fedops_provider.DigitalOceanProvider(auth)
  
  err := digo.CreateKeypair(sshKey)
  if err != nil {
    fmt.Println(err.Error())
    return d.Error
  }
  return d.Ok
}

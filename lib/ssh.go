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
  "fmt"
  "os"
  // 3rd Party
  "code.google.com/p/go.crypto/ssh"
  "code.google.com/p/go.crypto/ssh/terminal"
  // FedOps
)

func (d *Dispatcher) _bootstrap(vmID string, fedType uint) uint {
  ip := ""
  providerName := ""

  warehouses := d.Config.Warehouses
  for wIndex, _ := range warehouses {
    if warehouses[wIndex].WarehouseID == vmID {
      ip = warehouses[wIndex].IPV4
      providerName = warehouses[wIndex].Provider
      break
    }
  }

  trucks := d.Config.Trucks
  for tIndex, _ := range trucks {
    if trucks[tIndex].TruckID == vmID {
      ip = trucks[tIndex].IPV4
      providerName = trucks[tIndex].Provider
      break
    }
  }

  index := 0
  keys := d.Config.Keys
  for kIndex, _ := range keys {
    if keys[kIndex].ID[providerName] != "" {
      index = kIndex
      break
    }
  }

  key, err := ssh.ParsePrivateKey(d.Config.Keys[index].Keypair.PrivatePem)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  config := &ssh.ClientConfig{
    User: "root",
    Auth: []ssh.AuthMethod {
      ssh.PublicKeys(key),
    },
  }

  if ip == "" {
    fmt.Println("Could not find warehouse or truck with ID", vmID)
    return FedopsError
  }

  conn, err := ssh.Dial("tcp", ip + ":22", config)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer conn.Close()

  session, err := conn.NewSession()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer session.Close()

  // TODO :: Make this an external config file
  // Disable Password SSH login
  cmd := "sed --in-place=.bak 's/ChallengeResponseAuthentication\\ yes/ChallengeResponseAuthentication\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "sed --in-place=.bak 's/PasswordAuthentication\\ yes/PasswordAuthentication\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "sed --in-place=.bak 's/UsePAM\\ yes/UsePAM\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "systemctl restart sshd"
  // Generate a new server cert pair
  
  // TODO :: fedops user
  // Create a new fedops user, set sudoer settings
  // Install Docker, git and vim
  cmd += " && "
  cmd += "yum -y install docker git vim"
  cmd += " && "
  cmd += "systemctl start docker"
  cmd += " && "
  cmd += "systemctl enable docker"
  // Install Fedops
  cmd += " && "
  cmd += "docker build --no-cache=true --force-rm=true -t fedops github.com/wmiller848/Fedops"
  cmd += " && "
  // TODO :: Set up persistant container data
  if fedType == FedopsTypeTruck {
    cmd += "docker run --privileged -d -v=/opt/fedops:/opt/fedops/ fedops fedops-truck"
  } else if fedType == FedopsTypeWarehouse {
    cmd += "docker run --privileged -d -v=/opt/fedops:/opt/fedops/ fedops fedops-warehouse"
  }

  
  // fmt.Println("Running", cmd)
  err = session.Run(cmd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }  

  return FedopsOk
}


func (d *Dispatcher) SSH(vmID string) uint {
  d._ssh(vmID)
  // fmt.Println("SSH Session Ended")
  persisted := d.Unload()
  if persisted != true {
    return FedopsError
  } else {
    return FedopsOk
  }
}

func (d *Dispatcher) _ssh(vmID string) uint {

  promise := make(chan FedopsAction)
  go d.Refresh(promise)
  result := <- promise

  if result.Status == FedopsError {
    return FedopsError
  }

  ip := ""
  providerName := ""

  warehouses := d.Config.Warehouses
  for wIndex, _ := range warehouses {
    if warehouses[wIndex].WarehouseID == vmID {
      ip = warehouses[wIndex].IPV4
      providerName = warehouses[wIndex].Provider
      break
    }
  }

  trucks := d.Config.Trucks
  for tIndex, _ := range trucks {
    if trucks[tIndex].TruckID == vmID {
      ip = trucks[tIndex].IPV4
      providerName = trucks[tIndex].Provider
      break
    }
  }

  index := 0
  keys := d.Config.Keys
  for kIndex, _ := range keys {
    if keys[kIndex].ID[providerName] != "" {
      index = kIndex
      break
    }
  }

  key, err := ssh.ParsePrivateKey(d.Config.Keys[index].Keypair.PrivatePem)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  config := &ssh.ClientConfig{
    User: "root",
    Auth: []ssh.AuthMethod {
      ssh.PublicKeys(key),
    },
  }

  if ip == "" {
    fmt.Println("Could not find warehouse or truck with ID", vmID)
    return FedopsError
  }

  conn, err := ssh.Dial("tcp", ip + ":22", config)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer conn.Close()

  session, err := conn.NewSession()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer session.Close()

  fd := int(os.Stdin.Fd())
  oldState, err := terminal.MakeRaw(fd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer terminal.Restore(fd, oldState)

  session.Stdout = os.Stdout
  session.Stderr = os.Stderr
  session.Stdin = os.Stdin

  termWidth, termHeight, err := terminal.GetSize(fd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  // Set up terminal modes
  modes := ssh.TerminalModes{
    ssh.ECHO:          1,     // enable echoing
    ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
    ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
  }

  // Request pseudo terminal
  err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  err = session.Shell()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  // Will return the status of the last command run
  err = session.Wait()
  if err != nil {
    // fmt.Println(err.Error())
    return FedopsError
  }

  return FedopsOk
}
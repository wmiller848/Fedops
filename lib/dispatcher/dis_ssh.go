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
  "golang.org/x/crypto/ssh"
  "golang.org/x/crypto/ssh/terminal"
  "github.com/pkg/sftp"
  // FedOps
  "github.com/Fedops/lib/encryption"
  "github.com/Fedops/lib/engine"
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
  keys := d.Config.SSHKeys
  for kIndex, _ := range keys {
    if keys[kIndex].ID[providerName] != "" {
      index = kIndex
      break
    }
  }

  key, err := ssh.ParsePrivateKey(d.Config.SSHKeys[index].Keypair.PrivatePem)
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

  session, err := conn.NewSession()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  // session.Stdout = os.Stdout
  // session.Stderr = os.Stderr

  // TODO :: Make this an external config files 

  /////////////////
  // Disable Password SSH login and change port to 7575
  /////////////////
  cmd := "sed --in-place=.bak 's/ChallengeResponseAuthentication\\ yes/ChallengeResponseAuthentication\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "sed --in-place=.bak 's/PasswordAuthentication\\ yes/PasswordAuthentication\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "sed --in-place=.bak 's/UsePAM\\ yes/UsePAM\\ no/' /etc/ssh/sshd_config"
  cmd += " && "
  cmd += "sed --in-place=.bak 's/#Protocol\\ 2/Protocol\\ 2/' /etc/ssh/sshd_config"
  // cmd += " && "
  // cmd += "sed --in-place=.bak 's/#Port\\ 22/Port\\ 7575/' /etc/ssh/sshd_config"
  // cmd += " && "
  // cmd += "iptables -A INPUT -p tcp --dport 7575 -j ACCEPT"
  // cmd += " && "
  // cmd += "semanage port -a -t ssh_port_t -p tcp 7575"
  // Generate a new server cert pair
  
  /////////////////
  // TODO :: fedops user
  // Create a new fedops user, set sudoer settings
  /////////////////

  /////////////////
  // Install Docker, git, vim and sudo
  /////////////////
  cmd += " && "
  cmd += "yum -y install docker git"
  cmd += " && "
  cmd += "systemctl start docker"
  cmd += " && "
  cmd += "systemctl enable docker"

  /////////////////
  // Finally Restart SSHD
  /////////////////
  cmd += " && "
  cmd += "systemctl restart sshd"

  // fmt.Println("Running", cmd)
  err = session.Run(cmd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }


  session.Close()
  conn.Close()

  conn, err = ssh.Dial("tcp", ip + ":22", config)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  // session, err = conn.NewSession()
  // if err != nil {
  //   fmt.Println(err.Error())
  //   return FedopsError
  // }

  // session.Stdout = os.Stdout
  // session.Stderr = os.Stderr

  /////////////////
  // Install Fedops
  /////////////////
  // Write the config file
  keydata, err := fedops_encryption.GenerateRandomBytes(FedopsRemoteKeySize)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  pwd := "/opt/fedops"
  r := &fedops_runtime.Runtime{
    Cipherkey:      fedops_encryption.Encode(keydata),
    Config:         fedops_runtime.ClusterConfig{
      ClusterID: d.Config.ClusterID,
    },
    Version:        "0.0.1",
    PowerDirectory: pwd,
  }
  configData := r.UnloadToMemory()

  sftpClient, err := sftp.NewClient(conn)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  err = sftpClient.Mkdir(pwd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  configFile, err := sftpClient.Create(pwd + "/" + fedops_runtime.ConfigFileName)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  _, err = configFile.Write(configData)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  configFile.Close()

  keyFile, err := sftpClient.Create(pwd + "/" + fedops_runtime.KeyFileName)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  _, err = keyFile.Write(keydata)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  keyFile.Close()

  sftpClient.Close()
  conn.Close()

  conn, err = ssh.Dial("tcp", ip + ":22", config)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer conn.Close()

  session, err = conn.NewSession()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer session.Close()

  // Build fedops
  cmd = "docker build --no-cache=true --force-rm=true -t fedops " + FedopsRepo
  cmd += " && "
  // TODO :: Set up persistant data container instead of mounting a volume from the host
  if fedType == FedopsTypeTruck {
    cmd += "docker run --privileged -d -v=/opt/fedops:/opt/fedops/ fedops fedops-truck"
  } else if fedType == FedopsTypeWarehouse {
    cmd += "docker run --privileged -d -v=/opt/fedops:/opt/fedops/ fedops fedops-warehouse"
  }
  
  /////////////////
  // Execute the commands
  /////////////////
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
  keys := d.Config.SSHKeys
  for kIndex, _ := range keys {
    if keys[kIndex].ID[providerName] != "" {
      index = kIndex
      break
    }
  }

  key, err := ssh.ParsePrivateKey(d.Config.SSHKeys[index].Keypair.PrivatePem)
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
  userState, err := terminal.MakeRaw(fd)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  defer terminal.Restore(fd, userState)

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
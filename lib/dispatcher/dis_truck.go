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
  "time"
  // 3rd Party
  // FedOps
  "github.com/Fedops/lib/providers"
)

type Truck struct {
  fedops_provider.ProviderVM
  TruckID  string
  Containers []string
}

func (d *Dispatcher) CreateTruck(promise chan FedopsAction, provider, memSize, diskSize, numVcpus string) {
  // Cycle through all the provider tokens
  for name, token := range d.Config.Tokens {
    switch name {
      case fedops_provider.DigitalOceanName:
        auth := fedops_provider.DigitalOceanAuth{
          ApiKey: token.AccessToken,
        }
        provider := fedops_provider.DigitalOceanProvider(auth)
        status := d._createTruck(&provider)
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

func (d *Dispatcher) _createTruck(provider fedops_provider.Provider) uint {
  size, err := provider.GetDefaultSize()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  image, err := provider.GetDefaultImage()
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  // See if there is a key for this provider
  vmid, err := GenerateRandomHex(TruckIDSize)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  vm, err := provider.CreateVM(vmid, size, image, d.Config.SSHKeys)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  truck := new(Truck)
  truck.TruckID = vmid
  truck.ID = vm.ID
  truck.Provider = provider.Name()
  d.Config.Trucks = append(d.Config.Trucks, *truck)

  fmt.Printf("Initializing...")

  done := false
  for done == false {
    time.Sleep(FedopsPoolTime * time.Second)
    fmt.Printf(".")
    // fmt.Println("Refreshing...")
    promise := make(chan FedopsAction)
    go d.Refresh(promise)
    result := <- promise

    if result.Status == FedopsError {
      return FedopsError
    }

    trucks := d.Config.Trucks
    for tIndex, _ := range trucks {
      if trucks[tIndex].TruckID == truck.TruckID {
        if trucks[tIndex].Status == "up" {
          done = true
          break 
        }      
      }
    }
  }
  fmt.Printf("\r\n")

  // Give the machine a few seconds to boot
  time.Sleep(FedopsBootWaitTime * time.Second)

  fmt.Printf("Bootstrapping...")
  done = false
  go func() {
    for done == false {
      time.Sleep(FedopsPoolTime * time.Second)
      fmt.Printf(".")
    }
  }()

  d._bootstrap(truck.TruckID, FedopsTypeTruck)
  done = true
  fmt.Printf("\r\n")

  return FedopsOk
}

func (d *Dispatcher) DestroyTruck(promise chan FedopsAction, truckID string) {

  var truck Truck
  trucks := d.Config.Trucks
  found := false
  var tIndex int
  for tIndex = range trucks {
    if trucks[tIndex].TruckID == truckID {
      truck = trucks[tIndex]
      found = true
      break
    }
  }

  if !found {
    fmt.Println("Unable to locate truck with ID " + truckID)
    promise <- FedopsAction{
      Status: FedopsError,
    }
    return
  }

  token := d.Config.Tokens[truck.Provider]

  switch truck.Provider {
    case fedops_provider.DigitalOceanName:
      auth := fedops_provider.DigitalOceanAuth{
        ApiKey: token.AccessToken,
      }
      provider := fedops_provider.DigitalOceanProvider(auth)
      status := d._destroyTruck(&provider, truck)
      if status == FedopsError {
        promise <- FedopsAction{
          Status: FedopsError,
        }
        return
      }
  }

  d.Config.Trucks = append(d.Config.Trucks[:tIndex], d.Config.Trucks[tIndex+1:]...)

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

func (d *Dispatcher) _destroyTruck(provider fedops_provider.Provider, truck Truck) uint {
  vm := fedops_provider.ProviderVM{
    ID: truck.ID,
  }
  err := provider.DestroyVM(vm)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  return FedopsOk
}
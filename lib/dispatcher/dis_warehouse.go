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
  "github.com/Fedops/lib/encryption"
  "github.com/Fedops/lib/engine/warehouse"
)

// type Warehouse struct {
//   fedops_provider.ProviderVM
//   WarehouseID string
//   Containers    []string
// }

func (d *Dispatcher) CreateWarehouse(promise chan FedopsAction, providerName, memSize, diskSize, numVcpus string) {
  // Cycle through all the provider tokens
  for name, token := range d.Config.Tokens {
    switch name {
      case fedops_provider.DigitalOceanName:
        auth := fedops_provider.DigitalOceanAuth{
          ApiKey: token.AccessToken,
        }
        provider := fedops_provider.DigitalOceanProvider(auth)
        status := d._createWarehouse(&provider)
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

func (d *Dispatcher) _createWarehouse(provider fedops_provider.Provider) uint {
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
  vmid, err := fedops_encryption.GenerateRandomHex(WarehouseIDSize)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  vm, err := provider.CreateVM(vmid, size, image, d.Config.SSHKeys)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }

  warehouse := new(fedops_warehouse.Warehouse)
  warehouse.WarehouseID = vmid
  warehouse.ID = vm.ID
  warehouse.Provider = provider.Name()
  d.Config.Warehouses[vmid] = warehouse

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

    if d.Config.Warehouses[vmid].Status == "up" {
      done = true
    } 
  }
  fmt.Printf("\r\n")

  fmt.Printf("Bootstrapping...")
  done = false
  go func() {
    for done == false {
      time.Sleep(FedopsPoolTime * time.Second)
      fmt.Printf(".")
    }
  }()

  // Give the machine a few seconds to boot
  time.Sleep(FedopsBootWaitTime * time.Second)

  d._bootstrap(warehouse.WarehouseID, FedopsTypeWarehouse)
  done = true
  fmt.Printf("\r\n")

  return FedopsOk
}

func (d *Dispatcher) DestroyWarehouse(promise chan FedopsAction, warehouseID string) {

  warehouse, ok := d.Config.Warehouses[warehouseID]

  if !ok {
    fmt.Println("Unable to locate warehouse with ID " + warehouseID)
    promise <- FedopsAction{
      Status: FedopsError,
    }
    return
  }

  token := d.Config.Tokens[warehouse.Provider]

  switch warehouse.Provider {
    case fedops_provider.DigitalOceanName:
      auth := fedops_provider.DigitalOceanAuth{
        ApiKey: token.AccessToken,
      }
      provider := fedops_provider.DigitalOceanProvider(auth)
      status := d._destroyWarehouse(&provider, *warehouse)
      if status == FedopsError {
        promise <- FedopsAction{
          Status: FedopsError,
        }
        return
      }
  }

  d.Config.Warehouses[warehouseID] = nil
  delete(d.Config.Warehouses, warehouseID)

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

func (d *Dispatcher) _destroyWarehouse(provider fedops_provider.Provider, warehouse fedops_warehouse.Warehouse) uint {
  vm := fedops_provider.ProviderVM{
    ID: warehouse.ID,
  }
  err := provider.DestroyVM(vm)
  if err != nil {
    fmt.Println(err.Error())
    return FedopsError
  }
  return FedopsOk
}
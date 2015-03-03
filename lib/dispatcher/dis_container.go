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
  "fmt"
  //
  "golang.org/x/crypto/bcrypt"
  //
  "github.com/Fedops/lib/encryption"
  "github.com/Fedops/lib/engine/container"
  "github.com/Fedops/lib/engine/network"
)

// type Container struct {
//   ContainerID  string
//   Repo string
//   Warehouses []string
//   Trucks []string
// }

func (d *Dispatcher) CreateContainer(promise chan FedopsAction, repo string) {
  containerID, _ := fedops_encryption.GenerateRandomHex(ContainerIDSize)

  container := new(fedops_container.Container)
  container.Repo = repo
  container.ContainerID = containerID
  
  d.Config.Containers = append(d.Config.Containers, *container)

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

func (d *Dispatcher) DestroyContainer(promise chan FedopsAction, containerID string) {
  containers := d.Config.Containers
  found := false
  var cIndex int
  for cIndex = range containers {
    if containers[cIndex].ContainerID == containerID {
      found = true
      break 
    }
  }

  if !found {
    fmt.Println("Unable to locate container with ID " + containerID)
    promise <- FedopsAction{
      Status: FedopsError,
    }
    return
  }

  d.Config.Containers = append(d.Config.Containers[:cIndex], d.Config.Containers[cIndex+1:]...)

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


// Ship a container to the warehouse for continuous deployment
func (d *Dispatcher) _shipContainerToWarehouse(containerID, warehouseID string) uint {
  ip := ""
  warehouses := d.Config.Warehouses
  for wIndex, _ := range warehouses {
    if warehouses[wIndex].WarehouseID == warehouseID {
      ip = warehouses[wIndex].IPV4
      break
    }
  }

  if ip == "" {
    fmt.Println("Could not find warehouse with ID", warehouseID)
    return FedopsError
  }

  conn := d.OpenConnection(ip)
  defer conn.Close()

  req := fedops_network.FedopsRequest{
    Method: fedops_network.FedopsRequestCreate,
    Route: []byte("container"),
  }
  err := d.WriteToConn(conn, &req)
  if err != nil {
    fmt.Println(err.Error()) 
    return FedopsError
  }

  return FedopsOk
}

func (d *Dispatcher) ShipContainerToWarehouse(promise chan FedopsAction, containerID, warehouseID string) {

    success := d._shipContainerToWarehouse(containerID, warehouseID)
    persisted := d.Unload()
    if persisted != true || success == FedopsError {
      promise <- FedopsAction{
        Status: FedopsError,
      }
    }
    promise <- FedopsAction{
      Status: FedopsOk,
    }
}

// Ship a container to the warehouse for continuous deployment
func (d *Dispatcher) _shipContainerImageToTruck(containerID, truckID string) uint {
  ip := ""
  trucks := d.Config.Trucks
  for tIndex, _ := range trucks {
    if trucks[tIndex].TruckID == truckID {
      ip = trucks[tIndex].IPV4
      break
    }
  }

  if ip == "" {
    fmt.Println("Could not find truck with ID", truckID)
    return FedopsError
  }

  conn := d.OpenConnection(ip)
  defer conn.Close()

  auth, err := bcrypt.GenerateFromPassword([]byte(d.Config.ClusterID), AuthorizationCost)
  if err != nil {
    fmt.Println(err.Error()) 
    return FedopsError
  }

  req := fedops_network.FedopsRequest{
    Authorization: auth,
    Method: fedops_network.FedopsRequestCreate,
    Route: []byte("container"),
  }
  err = d.WriteToConn(conn, &req)
  if err != nil {
    fmt.Println(err.Error()) 
    return FedopsError
  }

  return FedopsOk
}

func (d *Dispatcher) ShipContainerImageToTruck(promise chan FedopsAction, containerID, truckID string) {

    success := d._shipContainerImageToTruck(containerID, truckID)
    persisted := d.Unload()
    if persisted != true || success == FedopsError {
      promise <- FedopsAction{
        Status: FedopsError,
      }
    }
    promise <- FedopsAction{
      Status: FedopsOk,
    }
}
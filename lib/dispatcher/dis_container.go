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
  
  d.Config.Containers[containerID] = container

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
  _, ok := d.Config.Containers[containerID]

  if !ok {
    fmt.Println("Unable to locate container with ID " + containerID)
    promise <- FedopsAction{
      Status: FedopsError,
    }
    return
  }

  d.Config.Containers[containerID] = nil

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

  container, ok := d.Config.Containers[containerID]
  if !ok {
    fmt.Println("Could not find container with ID", containerID)
    return FedopsError
  }

  warehouse, ok := d.Config.Warehouses[warehouseID]
  if !ok {
    fmt.Println("Could not find warehouse with ID", warehouseID)
    return FedopsError
  }

  auth, err := bcrypt.GenerateFromPassword([]byte(d.Config.ClusterID), AuthorizationCost)
  if err != nil {
    fmt.Println(err.Error()) 
    return FedopsError
  }

  // Find the state of this container
  // Does it already have any trucks
  if len(container.Trucks) > 0 {
    for key := range container.Trucks {
      truck := d.Config.Trucks[container.Trucks[key]]
      req := fedops_network.FedopsRequest{
        Authorization: auth,
        Method: fedops_network.FedopsRequestUpdate,
        Route: []byte("/container/" + containerID),
        Data: []byte("truck:" + container.Trucks[key]),
      }

      conn := d.OpenConnection(warehouse.IPV4)
      err = d.WriteToConn(conn, &req)
      if err != nil {
        fmt.Println(err.Error()) 
        return FedopsError
      }
      conn.Close()

      // Update the warehouse
      req = fedops_network.FedopsRequest{
        Authorization: auth,
        Method: fedops_network.FedopsRequestCreate,
        Route: []byte("/container/" + containerID),
        Data: []byte("warehouse:" + warehouseID),
      }

      conn = d.OpenConnection(truck.IPV4)
      err = d.WriteToConn(conn, &req)
      if err != nil {
        fmt.Println(err.Error()) 
        return FedopsError
      }
      conn.Close()
    }
    container.Warehouse = warehouseID
  } else {
    req := fedops_network.FedopsRequest{
      Authorization: auth,
      Method: fedops_network.FedopsRequestCreate,
      Route: []byte("/container/" + containerID),
    }

    conn := d.OpenConnection(warehouse.IPV4)
    defer conn.Close()

    err = d.WriteToConn(conn, &req)
    if err != nil {
      fmt.Println(err.Error()) 
      return FedopsError
    }

    container.Warehouse = warehouseID
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

  container, ok := d.Config.Containers[containerID]
  if !ok {
    fmt.Println("Could not find container with ID", containerID)
    return FedopsError
  }

  truck, ok := d.Config.Trucks[truckID]
  if !ok {
    fmt.Println("Could not find truck with ID", truckID)
    return FedopsError
  }

  auth, err := bcrypt.GenerateFromPassword([]byte(d.Config.ClusterID), AuthorizationCost)
  if err != nil {
    fmt.Println(err.Error()) 
    return FedopsError
  }

  // Find the state of this container
  // Does it already have a warehouse
  if container.Warehouse != "" {
    req := fedops_network.FedopsRequest{
      Authorization: auth,
      Method: fedops_network.FedopsRequestCreate,
      Route: []byte("/container/" + containerID),
      Data: []byte("warehouse:" + container.Warehouse),
    }

    conn := d.OpenConnection(truck.IPV4)
    err = d.WriteToConn(conn, &req)
    if err != nil {
      fmt.Println(err.Error()) 
      return FedopsError
    }
    conn.Close()

    // Update the warehouse
    req = fedops_network.FedopsRequest{
      Authorization: auth,
      Method: fedops_network.FedopsRequestCreate,
      Route: []byte("/container/" + containerID),
      Data: []byte("truck:" + truckID),
    }

    container.Trucks = append(container.Trucks, truckID)

    warehouse := d.Config.Warehouses[container.Warehouse]
    conn = d.OpenConnection(warehouse.IPV4)
    defer conn.Close()
    err = d.WriteToConn(conn, &req)
    if err != nil {
      fmt.Println(err.Error()) 
      return FedopsError
    }
    
  } else {
    req := fedops_network.FedopsRequest{
      Authorization: auth,
      Method: fedops_network.FedopsRequestCreate,
      Route: []byte("/container/" + containerID),
    }

    conn := d.OpenConnection(truck.IPV4)
    defer conn.Close()

    err = d.WriteToConn(conn, &req)
    if err != nil {
      fmt.Println(err.Error()) 
      return FedopsError
    }

    container.Trucks = append(container.Trucks, truckID)
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
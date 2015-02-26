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
  "github.com/Fedops/lib/encryption"
  "github.com/Fedops/lib/engine/container"
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
func (d *Dispatcher) ShipContainerToWarehouse(promise chan FedopsAction, containerID string) {
  d.OpenConnection(containerID)

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

// Ship a containers image to a truck for execution
func (d *Dispatcher) ShipContainerImageToTruck(promise chan FedopsAction, containerID string) {

}
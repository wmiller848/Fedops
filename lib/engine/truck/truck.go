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

package fedops_truck

import (
	//
	"bytes"
	"errors"
	"fmt"
	"os"
	//
	"github.com/wmiller848/Fedops/lib/engine"
	"github.com/wmiller848/Fedops/lib/engine/container"
	"github.com/wmiller848/Fedops/lib/engine/network"
	"github.com/wmiller848/Fedops/lib/providers"
)

type Truck struct {
	fedops_provider.ProviderVM
	TruckID    string
	Containers []string
}

type TruckDaemon struct {
	fedops_runtime.Runtime
}

func CreateDaemon() *TruckDaemon {
	pwd := os.Getenv("PWD")

	truckDaemon := TruckDaemon{}
	// Set up the default runtime
	err := truckDaemon.Configure(pwd)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	// Set up the routes for network calls
	err = truckDaemon.AddRoute(fedops_network.FedopsRequestInfo, "^/containers$", truckDaemon.ListContainers)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = truckDaemon.AddRoute(fedops_network.FedopsRequestCreate, "^/container/[A-Za-z0-9]+$", truckDaemon.ShipContainer)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = truckDaemon.AddRoute(fedops_network.FedopsRequestUpdate, "^/container/[A-Za-z0-9]+$", truckDaemon.UpdateContainer)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = truckDaemon.AddRoute(fedops_network.FedopsRequestDestroy, "^/container/[A-Za-z0-9]+$", truckDaemon.UnshipContainer)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = truckDaemon.AddRoute(fedops_network.FedopsRequestCreate, "^/container/[A-Za-z0-9]+/[A-Za-z0-9]+$", truckDaemon.ShipContainerImage)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &truckDaemon
}

func (d *TruckDaemon) ListContainers(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("ListContainers", string(req.Data), args)
	return nil
}

func (d *TruckDaemon) ShipContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	var containerID, warehouseID string
	args := bytes.Split(req.Route, []byte("/"))
	if len(args) >= 3 {
		containerID = string(args[2])
	}
	dataArgs := bytes.Split(req.Data, []byte(":"))
	if len(dataArgs) >= 2 {
		warehouseID = string(dataArgs[1])
	}

	fmt.Println(args, dataArgs)

	if containerID == "" {
		return errors.New("Bad ContainerID")
	}
	fmt.Println("ShipContainer", containerID, warehouseID)

	d.Config.Containers[containerID] = fedops_container.Container{
		ContainerID: containerID,
		Warehouse:   warehouseID,
	}
	return nil
}

func (d *TruckDaemon) UpdateContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	var containerID, warehouseID string
	args := bytes.Split(req.Route, []byte("/"))
	if len(args) >= 3 {
		containerID = string(args[2])
	}
	dataArgs := bytes.Split(req.Data, []byte(":"))
	if len(dataArgs) >= 2 {
		warehouseID = string(dataArgs[1])
	}

	fmt.Println(args, dataArgs)

	if containerID == "" {
		return errors.New("Bad ContainerID")
	}
	fmt.Println("UpdateContainer", containerID, warehouseID)

	container := d.Config.Containers[containerID]
	container.Warehouse = warehouseID
	return nil
}

func (d *TruckDaemon) UnshipContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	var containerID string
	args := bytes.Split(req.Route, []byte("/"))
	if len(args) >= 3 {
		containerID = string(args[2])
	}

	fmt.Println(args)

	if containerID == "" {
		return errors.New("Bad ContainerID")
	}
	fmt.Println("UnshipContainer", containerID)

	delete(d.Config.Containers, containerID)
	return nil
}

func (d *TruckDaemon) ShipContainerImage(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("ShipContainerImage", string(req.Data), args)
	return nil
}

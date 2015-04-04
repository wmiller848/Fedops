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

package fedops_warehouse

import (
	//
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
	//
	"github.com/wmiller848/Fedops/lib/engine"
	"github.com/wmiller848/Fedops/lib/engine/network"
	"github.com/wmiller848/Fedops/lib/providers"
)

type Warehouse struct {
	fedops_provider.ProviderVM
	WarehouseID string
	Containers  []string
}

type WarehouseDaemon struct {
	fedops_runtime.Runtime
}

func CreateDaemon() *WarehouseDaemon {
	pwd := os.Getenv("PWD")

	warehouseDaemon := WarehouseDaemon{}
	// Set up the default runtime
	warehouseDaemon.Configure(pwd)
	// Set up the routes for network calls
	err := warehouseDaemon.AddRoute(fedops_network.FedopsRequestInfo, "^/containers$", warehouseDaemon.ListContainers)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = warehouseDaemon.AddRoute(fedops_network.FedopsRequestCreate, "^/container/[A-Za-z0-9]+$", warehouseDaemon.PackageContainer)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = warehouseDaemon.AddRoute(fedops_network.FedopsRequestUpdate, "^/container/[A-Za-z0-9]+$", warehouseDaemon.UpdateContainer)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = warehouseDaemon.AddRoute(fedops_network.FedopsRequestDestroy, "^/container/[A-Za-z0-9]+$", warehouseDaemon.UnpackageContainer)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = warehouseDaemon.AddRoute(fedops_network.FedopsRequestCreate, "^/container/[A-Za-z0-9]+/[A-Za-z0-9]+$", warehouseDaemon.PackageContainerImage)
	if err != nil {
		fmt.Println(err.Error())
	}
	return &warehouseDaemon
}

func (d *WarehouseDaemon) ListContainers(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("ListContainers", string(req.Data), args)
	return nil
}

func (d *WarehouseDaemon) PackageContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	var containerID, truckID string
	args := bytes.Split(req.Route, []byte("/"))
	if len(args) >= 3 {
		containerID = string(args[2])
	}
	dataArgs := bytes.Split(req.Data, []byte(":"))
	if len(dataArgs) >= 2 {
		truckID = string(dataArgs[1])
	}
	fmt.Println("PackageContainer", containerID, truckID)
	event := fedops_runtime.FedopsEvent{
		ID:         containerID + ":" + truckID,
		Handle:     d.PollSourceControll,
		Persistant: true,
		Time:       time.Now(),
	}
	fmt.Println(event)
	d.Events = append(d.Events, event)

	return nil
}

func (d *WarehouseDaemon) UpdateContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("UpdateContainer", string(req.Data), args)
	return nil
}

func (d *WarehouseDaemon) UnpackageContainer(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("UnpackageContainer", string(req.Data), args)
	return nil
}

func (d *WarehouseDaemon) PackageContainerImage(req *fedops_network.FedopsRequest, res *fedops_network.FedopsResponse) error {
	args := bytes.Split(req.Route, []byte("/"))
	fmt.Println("PackageContainerImage", string(req.Data), args)
	return nil
}

func (d *WarehouseDaemon) PollSourceControll(event *fedops_runtime.FedopsEvent) {
	fmt.Println("PollSourceControll", event)
	idArgs := strings.Split(event.ID, ":")
	containerID := idArgs[0]
	// truckID := idArgs[0]

	container := d.Config.Containers[containerID]

	cmd := exec.Command("git clone")

	fmt.Println(container, cmd)
}

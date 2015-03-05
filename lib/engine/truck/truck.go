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
  "fmt"
  "os"
  //
  "github.com/Fedops/lib/providers"
  "github.com/Fedops/lib/engine"
  "github.com/Fedops/lib/engine/network"
)

type Truck struct {
  fedops_provider.ProviderVM
  TruckID  string
  Containers []string
}

type TruckDaemon struct {
  fedops_runtime.Runtime
}

func CreateDaemon() *TruckDaemon{
  pwd := os.Getenv("PWD")

  truckDaemon := TruckDaemon{}
  // Set up the default runtime
  truckDaemon.Configure(pwd)
  // Set up the routes for network calls
  truckDaemon.AddRoute("^/container$", truckDaemon.ShipContainer)
  truckDaemon.AddRoute("^/container/[A-Za-z0-9]+$", truckDaemon.UnshipContainer)
  return &truckDaemon
}

func (d *TruckDaemon) ShipContainer(req fedops_network.FedopsRequest) error {
  fmt.Println("SHIP", req)
  return nil
}

func (d *TruckDaemon) UnshipContainer(req fedops_network.FedopsRequest) error {
  fmt.Println("UNSHIP", req)
  return nil
}
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
  "regexp"
  "net"
  "crypto/tls"
  "encoding/gob"
  //
  "golang.org/x/crypto/bcrypt"
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
  Routes []fedops_network.FedopsRoute
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

///////
// TODO :: move into runtime
///////
func (d *TruckDaemon) AddRoute(route string, handle fedops_network.HandleRoute) error {
  rgx, err := regexp.Compile(route)
  if err != nil {
    return err
  }
  fedRoute := fedops_network.FedopsRoute{
    Route: rgx,
    Handle: handle,
  }
  d.Routes = append(d.Routes, fedRoute)
  return nil
}

// Handles incoming requests.
func (d *TruckDaemon) HandleConnection(conn net.Conn) {
  // Make a buffer to hold incoming data.
  // buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  // reqLen, err := conn.Read(buf)
  // if err != nil {
  //   fmt.Println("Error reading:", err.Error())
  //   return
  // }
  // if reqLen > 0 {
  //   fmt.Println(buf)
  // }
  // Send a response back to person contacting us.
  // conn.Write([]byte("Message received"))
  defer conn.Close()
  dec := gob.NewDecoder(conn)
  var req fedops_network.FedopsRequest
  err := dec.Decode(&req)
  if err != nil {
    fmt.Println(err.Error())
    return
  }
  err = bcrypt.CompareHashAndPassword(req.Authorization, []byte(d.Config.ClusterID))
  if err != nil {
    fmt.Println("Authorization not accepted", err.Error())
    return
  } else {
    fmt.Println("Authorization accepted")
    fmt.Println("Method", req.Method)
    fmt.Println("Route", string(req.Route))

    for i := range d.Routes {
      if d.Routes[i].Route.Match(req.Route) {
        err = d.Routes[i].Handle(req)
        if err != nil {
          fmt.Println(err.Error())
        }
        break
      }
    }
  }
  conn.Write([]byte("ok"))
}

func (d *TruckDaemon) Listen() {
  // config := &ssh.ServerConfig{}
  // private, err := ssh.ParsePrivateKey(d.Config.)
  // if err != nil {
  //   log.Fatal("Failed to parse private key")
  // }
 
  // config.AddHostKey(private)

  fed_cert := d.Config.Cert
  // cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
  cert, err := tls.X509KeyPair(fed_cert.CertificatePem, fed_cert.PrivatePem)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  config := tls.Config{Certificates: []tls.Certificate{cert}}
  listener, err := tls.Listen("tcp", ":13371", &config)
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  for {
      conn, err := listener.Accept()
      if err != nil {
        fmt.Println(err.Error())
        break
      }
      fmt.Println(conn.RemoteAddr(), "Connected")
      go d.HandleConnection(conn)
  }
}
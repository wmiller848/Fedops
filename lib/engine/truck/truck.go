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
  //
  "github.com/Fedops/lib/providers"
  "github.com/Fedops/lib/engine"
)

type Truck struct {
  fedops_provider.ProviderVM
  TruckID  string
  Containers []string
}

type TruckDaemon struct {
  fedops_runtime.Runtime
  Muxer regexp.Regexp
}

func CreateDaemon() *TruckDaemon{
  pwd := os.Getenv("PWD")

  truckDaemon := TruckDaemon{}
  // Set up the default runtime
  truckDaemon.Configure(pwd)
  return &truckDaemon
}

// Handles incoming requests.
func handleConnection(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  reqLen, err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  fmt.Println(reqLen, buf)
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received"))
  // Close the connection when you're done with it.
  conn.Close()
}

func (d *TruckDaemon) Listen() {
  fed_certs := d.Config.Certs

  // cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
  cert, err := tls.X509KeyPair(fed_certs[0].CertificatePem, fed_certs[0].PrivatePem)
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
      go handleConnection(conn)
  }

}
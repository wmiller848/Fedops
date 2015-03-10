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
  "errors"
  "crypto/tls"
  "crypto/x509"
  "encoding/gob"
  //
  "github.com/Fedops/lib/engine/network"
)

func (d *Dispatcher) OpenConnection(vmIP string) *tls.Conn {
  // server cert is self signed -> server_cert == ca_cert
  certPool := x509.NewCertPool()

  fed_certs := d.Config.Certs
  certPool.AppendCertsFromPEM(fed_certs[0].CertificatePem)

  config := tls.Config{
    RootCAs: certPool,
    PreferServerCipherSuites: true,
    SessionTicketsDisabled: true,
    CipherSuites: []uint16{
      tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    },
    CurvePreferences: []tls.CurveID{tls.CurveP521},
    MinVersion: tls.VersionTLS12,
    MaxVersion: tls.VersionTLS12,
  }

  conn, err := tls.Dial("tcp", vmIP + ":13371", &config)
  if err != nil {
    fmt.Println("client: dial:", err.Error())
    return nil
  }
  // defer conn.Close()
  return conn
}

func (d *Dispatcher) WriteToConn(conn *tls.Conn, req *fedops_network.FedopsRequest) error {
  enc := gob.NewEncoder(conn)
  err := enc.Encode(req)
  if err != nil {
    return err
  }

  dec := gob.NewDecoder(conn)
  var res fedops_network.FedopsResponse
  err = dec.Decode(&res)
  if err != nil {
    return err
  }

  if !res.Success {
    return errors.New("Remote returned " + string(res.Error))
  }

  return nil
}

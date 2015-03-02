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
  "crypto/tls"
  "crypto/x509"
  //
)

func (d *Dispatcher) OpenConnection(vmID string) *tls.Conn {
  // server cert is self signed -> server_cert == ca_cert
  certPool := x509.NewCertPool()

  // severCert, err := ioutil.ReadFile("./cert.pem")
  // if err != nil {
  //     log.Fatal("Could not load server certificate!")
  // }

  fed_certs := d.Config.Certs
  certPool.AppendCertsFromPEM(fed_certs[0].CertificatePem)

  config := tls.Config{RootCAs: certPool}

  conn, err := tls.Dial("tcp", "127.0.0.1:13371", &config)
  if err != nil {
      fmt.Println("client: dial:", err.Error())
      return nil
  }
  // defer conn.Close()
  return conn
}

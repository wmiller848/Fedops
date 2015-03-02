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

package fedops_encryption

import (
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/rand"
  // _ "crypto/rsa"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/pem"
  "fmt"
  "log"
  "math/big"
  "net"
  "os"
  "time"
)


type Cert_Config struct {
  IP string
}

func GenerateCert(certConfig Cert_Config) Cert {
  cert := Cert{}
  cert.Generate(certConfig)
  return cert
}

type Cert struct {
  // host string
  // ValidFrom string
  // ValidTo string
  // CA bool
  // KeySize int
  CertificatePem []byte
  PublicPem  []byte
  PrivatePem []byte
}

func (c *Cert) Generate(certConfig Cert_Config) {
  // TODO :: use the fedops keypair type
  // fmt.Println("Creating EC Key")
  priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
  if err != nil {
    fmt.Println(err)
    return
  }

  priv_der, err := x509.MarshalECPrivateKey(priv)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
    return
  }
  priv_blk := pem.Block{
    // Type:    "RSA PRIVATE KEY",
    Type:    "EC PRIVATE KEY",
    Headers: nil,
    Bytes:   priv_der,
  }

  // Public Key generation
  pub := priv.PublicKey
  pub_der, err := x509.MarshalPKIXPublicKey(&pub)
  if err != nil {
    fmt.Println("Failed to get der format for PublicKey.", err)
    return
  }

  pub_blk := pem.Block{
    Type:    "PUBLIC KEY",
    Headers: nil,
    Bytes:   pub_der,
  }

  notBefore := time.Now()
  var validFor time.Duration
  validFor = 365 * 24 * time.Hour
  notAfter := notBefore.Add(validFor)

  serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
  serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
  if err != nil {
    log.Fatalf("failed to generate serial number: %s", err)
  }

  template := x509.Certificate{
    DNSNames: []string{"*"},
    IPAddresses: []net.IP{net.ParseIP(certConfig.IP)},
    SerialNumber: serialNumber,
    Subject: pkix.Name{
      Organization: []string{"Fedops Daemon Certificate"},
    },
    NotBefore: notBefore,
    NotAfter:  notAfter,
    KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
    ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
    IsCA: true,
  }

  derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
  if err != nil {
    log.Fatalf("Failed to create certificate: %s", err)
  }

  cert_blk := pem.Block{
    Type: "CERTIFICATE",
    Bytes: derBytes,
  }

  //////////////////////////
  // Write to memory
  //////////////////////////

  c.CertificatePem = pem.EncodeToMemory(&cert_blk)
  c.PrivatePem = pem.EncodeToMemory(&priv_blk)
  c.PublicPem = pem.EncodeToMemory(&pub_blk)

  // //////////////////////////
  // // Write to disk
  // //////////////////////////
  // os.Mkdir("./.security", 0777)
  // fmt.Println("Writing Cert")
  // certOut, err := os.Create("./.security/cert.pem")
  // if err != nil {
  //   log.Fatalf("failed to open cert.pem for writing: %s", err)
  // }
  // pem.Encode(certOut, &cert_blk)
  // certOut.Close()
  // log.Print("written cert.pem\n")

  // fmt.Println("Writing Key")
  // keyOut, err := os.OpenFile("./.security/key.pem", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
  // if err != nil {
  //   log.Print("failed to open key.pem for writing:", err)
  //   return
  // }

  // pem.Encode(keyOut, &priv_blk)
  // keyOut.Close()
  // log.Print("written key.pem\n")
}
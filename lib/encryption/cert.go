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

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

package fedops_encryption

import (
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/rand"
  _ "crypto/rsa"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/pem"
  "fmt"
  "log"
  "math/big"
  "net"
  "os"
  _ "strings"
  "time"
)


type Cert_Config struct {

}

func GenerateCert(certConfig Cert_Config) Cert {
  cert := Cert{}
  cert.Generate()
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

func (c *Cert) Generate() {
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
    DNSNames: []string{"localhost"},
    IPAddresses: []net.IP{net.IP("127.0.0.1")},
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

// var (
//   host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
//   validFrom  = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
//   validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
//   isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
//   rsaBits    = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
//   ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256, P384, P521")
// )

// func publicKey(priv interface{}) interface{} {
//   switch k := priv.(type) {
//   case *rsa.PrivateKey:
//     return &k.PublicKey
//   case *ecdsa.PrivateKey:
//     return &k.PublicKey
//   default:
//     return nil
//   }
// }

// func pemBlockForKey(priv interface{}) *pem.Block {
//   switch k := priv.(type) {
//   case *rsa.PrivateKey:
//     return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
//   case *ecdsa.PrivateKey:
//     b, err := x509.MarshalECPrivateKey(k)
//     if err != nil {
//       fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
//       os.Exit(2)
//     }
//     return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
//   default:
//     return nil
//   }
// }

// func Generate() {
//   flag.Parse()

//   if len(*host) == 0 {
//     log.Fatalf("Missing required --host parameter")
//   }

//   var priv interface{}
//   var err error
//   switch *ecdsaCurve {
//   case "":
//     priv, err = rsa.GenerateKey(rand.Reader, *rsaBits)
//   case "P224":
//     priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
//   case "P256":
//     priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
//   case "P384":
//     priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
//   case "P521":
//     priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
//   default:
//     fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", *ecdsaCurve)
//     os.Exit(1)
//   }
//   if err != nil {
//     log.Fatalf("failed to generate private key: %s", err)
//   }

//   var notBefore time.Time
//   if len(*validFrom) == 0 {
//     notBefore = time.Now()
//   } else {
//     notBefore, err = time.Parse("Jan 2 15:04:05 2006", *validFrom)
//     if err != nil {
//       fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
//       os.Exit(1)
//     }
//   }

//   notAfter := notBefore.Add(*validFor)

//   serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
//   serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
//   if err != nil {
//     log.Fatalf("failed to generate serial number: %s", err)
//   }

//   template := x509.Certificate{
//     SerialNumber: serialNumber,
//     Subject: pkix.Name{
//       Organization: []string{"Acme Co"},
//     },
//     NotBefore: notBefore,
//     NotAfter:  notAfter,

//     KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
//     ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
//     BasicConstraintsValid: true,
//   }

//   hosts := strings.Split(*host, ",")
//   for _, h := range hosts {
//     if ip := net.ParseIP(h); ip != nil {
//       template.IPAddresses = append(template.IPAddresses, ip)
//     } else {
//       template.DNSNames = append(template.DNSNames, h)
//     }
//   }

//   if *isCA {
//     template.IsCA = true
//     template.KeyUsage |= x509.KeyUsageCertSign
//   }

//   derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
//   if err != nil {
//     log.Fatalf("Failed to create certificate: %s", err)
//   }

//   certOut, err := os.Create("cert.pem")
//   if err != nil {
//     log.Fatalf("failed to open cert.pem for writing: %s", err)
//   }
//   pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
//   certOut.Close()
//   log.Print("written cert.pem\n")

//   keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
//   if err != nil {
//     log.Print("failed to open key.pem for writing:", err)
//     return
//   }
//   pem.Encode(keyOut, pemBlockForKey(priv))
//   keyOut.Close()
//   log.Print("written key.pem\n")
// }
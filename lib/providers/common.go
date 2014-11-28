package fedops_provider

import (
  "fmt"
  "bytes"
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
  _"encoding/base64"
  "code.google.com/p/go.crypto/ssh"
)

type ProviderKeypair struct {
  ID string
}

type ProviderSize struct {
  ID string
  Memory float64
  Vcpus float64
  Disk float64
  Bandwidth float64
  Price float64
}

type ProviderImage struct {
  ID string
  Distribution string
  Version string
}

type ProviderVM struct {
  ID string
}

//
//
type Provider interface {
  Name() string

  CreateKeypair(string, Keypair) (ProviderKeypair, error)

  ListSize() ([]ProviderSize, error)
  GetDefaultSize() (ProviderSize, error)

  ListImage() ([]ProviderImage, error)
  GetDefaultImage() (ProviderImage, error)

  ListVM() ([]ProviderVM, error)
  CreateVM(string, ProviderSize, ProviderImage, ProviderKeypair) (ProviderVM, error)
  SnapShotVM(ProviderVM) (ProviderImage, error)
}

type SSH_Config struct {
  Keysize int
}

func GenerateKeypair(sshKeyConfig SSH_Config) (Keypair) {
  sshkey := Keypair{Keysize: sshKeyConfig.Keysize}
  sshkey.Generate()
  return sshkey
}

//
//
//
type Keypair struct {
  Keysize int
  PublicPem []byte
  PrivatePem []byte
  PublicSSH []byte
}

func (k *Keypair) Generate() {
  priv, err := rsa.GenerateKey(rand.Reader, k.Keysize);
  if err != nil {
    fmt.Println(err);
    return;
  }
  err = priv.Validate();
  if err != nil {
    fmt.Println("Validation failed.", err);
  }

  // Get der format. priv_der []byte
  priv_der := x509.MarshalPKCS1PrivateKey(priv);

  // pem.Block
  // blk pem.Block
  priv_blk := pem.Block {
    Type: "RSA PRIVATE KEY",
    Headers: nil,
    Bytes: priv_der,
  };

  // Resultant private key in PEM format.
  // priv_pem string
  k.PrivatePem = bytes.Trim(pem.EncodeToMemory(&priv_blk), "\n")

  // Public Key generation
  pub := priv.PublicKey;
  pub_der, err := x509.MarshalPKIXPublicKey(&pub);
  if err != nil {
    fmt.Println("Failed to get der format for PublicKey.", err);
    return
  }

  pub_blk := pem.Block {
    Type: "PUBLIC KEY",
    Headers: nil,
    Bytes: pub_der,
  }
  k.PublicPem = bytes.Trim(pem.EncodeToMemory(&pub_blk), "\n")

  pubssh, err := ssh.NewPublicKey(&pub)
  if err != nil {
      fmt.Println("Failed to get ssh format for PublicKey.", err)
      return
  }
  k.PublicSSH = bytes.Trim(ssh.MarshalAuthorizedKey(pubssh), "\n")
}

func (k *Keypair) ToArray() []byte {
  return append(k.PrivatePem, k.PublicPem...)
}

func (k *Keypair) ToString() string {
  return string(k.PrivatePem) + string(k.PublicPem)
}
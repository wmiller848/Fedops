package fedops_provider

import (
  "fmt"
  "crypto/rsa"
  "crypto/rand"
  "crypto/x509"
  "encoding/pem"
)

//
//
type Provider interface {
  CreateKeypair()
  CreateImage()
  CreateVM()
}

type SSH_Config struct {
  Keysize int
}

//
//
//
type Keypair struct {
  Keysize int
  PublicPem []byte
  PrivatePem []byte
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
  k.PrivatePem = pem.EncodeToMemory(&priv_blk)

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
  k.PublicPem = pem.EncodeToMemory(&pub_blk)
}

func (k *Keypair) ToArray() []byte {
  return append(k.PrivatePem, k.PublicPem...)
}

func (k *Keypair) ToString() string {
  return string(k.PrivatePem) + string(k.PublicPem)
}
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

package fedops_provider

import (
	_ "bytes"
	"code.google.com/p/go.crypto/ssh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	_ "encoding/base64"
	"encoding/pem"
	"fmt"
)

// The Provider interface describes the functionality required for fedops
// to use a given provider, this includes, among other things, creating a new ssh keypair,
// listing image types, listing virtual machines, and reating virtual machines
type Provider interface {
	Name() string

	CreateKeypair(string, Keypair) (ProviderKeypair, error)

	ListSize() (ProviderSize, error)
	ListSizes() ([]ProviderSize, error)
	GetDefaultSize() (ProviderSize, error)

	ListImage() (ProviderImage, error)
	ListImages() ([]ProviderImage, error)
	GetDefaultImage() (ProviderImage, error)

	ListVM() (ProviderVM, error)
	ListVMs() ([]ProviderVM, error)
	CreateVM(string, ProviderSize, ProviderImage, []ProviderKeypair) (ProviderVM, error)

  DestroyVM(ProviderVM) error

	SnapShotVM(ProviderVM) (ProviderImage, error)
}

// The ProviderKeypair describes an id map that links a given provider id
// and the keypair together
type ProviderKeypair struct {
	ID      map[string]string
	Keypair Keypair
}

// The ProviderSize describes an type of size that may be used
// when creating a new virtual machine
type ProviderSize struct {
	ID        map[string]string
	Memory    float64
	Vcpus     float64
	Disk      float64
	Bandwidth float64
	Price     float64
}

type ProviderImage struct {
	ID           map[string]string
	Distribution string
	Version      string
}

type ProviderVM struct {
	ID       map[string]string
	IPV4     string
	IPV6     string
	Provider string
	Status   string
}

type SSH_Config struct {
	Keysize int
}

func GenerateKeypair(sshKeyConfig SSH_Config) Keypair {
	sshkey := Keypair{Keysize: sshKeyConfig.Keysize}
	sshkey.Generate()
	return sshkey
}

//
//
//
type Keypair struct {
	Keysize    int
	PublicPem  []byte
	PrivatePem []byte
	PublicSSH  []byte
}

func (k *Keypair) Generate() {
	priv, err := rsa.GenerateKey(rand.Reader, k.Keysize)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = priv.Validate()
	if err != nil {
		fmt.Println("Validation failed.", err)
	}

	// Get der format. priv_der []byte
	priv_der := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	// blk pem.Block
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}

	// Resultant private key in PEM format.
	// priv_pem string
	//k.PrivatePem = bytes.Trim(pem.EncodeToMemory(&priv_blk), "\n")
	k.PrivatePem = pem.EncodeToMemory(&priv_blk)

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
	//k.PublicPem = bytes.Trim(pem.EncodeToMemory(&pub_blk), "\n")
	k.PublicPem = pem.EncodeToMemory(&pub_blk)

	pubssh, err := ssh.NewPublicKey(&pub)
	if err != nil {
		fmt.Println("Failed to get ssh format for PublicKey.", err)
		return
	}
	//k.PublicSSH = bytes.Trim(ssh.MarshalAuthorizedKey(pubssh), "\n")
	k.PublicSSH = ssh.MarshalAuthorizedKey(pubssh)
}

func (k *Keypair) ToArray() []byte {
	return append(k.PrivatePem, k.PublicPem...)
}

func (k *Keypair) ToString() string {
	return string(k.PrivatePem) + string(k.PublicPem)
}

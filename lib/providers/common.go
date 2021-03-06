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
	_ "fmt"

	"github.com/wmiller848/Fedops/lib/encryption"
)

// The Provider interface describes the functionality required for fedops
// to use a given provider, this includes, among other things, creating a new ssh keypair,
// listing image types, listing virtual machines, and reating virtual machines
type Provider interface {
	Name() string

	CreateKeypair(string, fedops_encryption.Keypair) (ProviderKeypair, error)

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

// The ProviderCert describes an id map that links a given provider id
// and the cert together
type ProviderCerts struct {
	ID   map[string]string
	Cert fedops_encryption.Cert
}

// The ProviderKeypair describes an id map that links a given provider id
// and the keypair together
type ProviderKeypair struct {
	ID      map[string]string
	Keypair fedops_encryption.Keypair
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

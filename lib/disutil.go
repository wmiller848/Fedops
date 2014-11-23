package fedops

import (
  // Standard
  "crypto/rand"
  "encoding/base64"
  // 3rd Party
  // FedOps
  "github.com/FedOps/lib/providers"
)

const (
  DigitalOcean uint = 0
  AWS uint = 1
  GoogleCloud uint = 2
  MicrosoftAzure uint = 3
  OpenStack uint = 4
)

type ProviderTokens struct {
  AccessToken string
  SecurityToken string
}

// SSH Keypairs for clusters by provider
type ProviderKeys struct {
  DigitalOcean map[string]fedops_provider.Keypair
  AWS map[string]fedops_provider.Keypair
  GoogleCloud map[string]fedops_provider.Keypair
  MicrosoftAzure map[string]fedops_provider.Keypair
  OpenStack map[string]fedops_provider.Keypair
}

type VM struct {
  Provider string
  Role uint
  IP string
  Aliases []string
}

// GenerateRandomBytes returns securely generated random bytes
func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
      return nil, err
    }
    return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded, securely generated random string
func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.URLEncoding.EncodeToString(b), err
}
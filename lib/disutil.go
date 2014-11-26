package fedops

import (
  // Standard
  "crypto/rand"
  "encoding/base64"
  // 3rd Party
  // FedOps
  _"github.com/FedOps/lib/providers"
)

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

func decrypt(bytz []byte) []byte {
  return bytz
}

func encrypt(bytz []byte) []byte {
  return bytz
}

func decode(bytz []byte) []byte {
  return bytz
}

func encode(bytz []byte) []byte {
  return bytz
}
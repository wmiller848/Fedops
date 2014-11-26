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
    return base64.StdEncoding.EncodeToString(b), err
}

func decrypt(bytz []byte) []byte {
  return bytz
}

func encrypt(bytz []byte) []byte {
  return bytz
}

func decode(bytz []byte) []byte {
  _r, _ := base64.StdEncoding.DecodeString(string(bytz))
  return _r
}

func encode(bytz []byte) []byte {
  _r := []byte(base64.StdEncoding.EncodeToString(bytz))
  return _r
}
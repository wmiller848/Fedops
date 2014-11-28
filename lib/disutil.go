package fedops

import (
  // Standard
  _"crypto/md5"
  "crypto/sha256"
  "crypto/rand"
  "crypto/aes"
  "crypto/cipher"
  "encoding/base64"
  "io"
  // 3rd Party
  // FedOps
)

func GenerateRandomBytes(n int) ([]byte, error) {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
      return nil, err
    }
    return b, nil
}

func GenerateRandomString(s int) (string, error) {
    b, err := GenerateRandomBytes(s)
    return base64.StdEncoding.EncodeToString(b), err
}

func Encode(bytz []byte) ([]byte) {
  return []byte(base64.StdEncoding.EncodeToString(bytz))  
}

func Decode(bytz []byte) ([]byte, error) {
  return base64.StdEncoding.DecodeString(string(bytz))
}

func Encrypt(key, bytz []byte) ([]byte, error) {
  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }
  b := Encode(bytz)
  ciphertext := make([]byte, aes.BlockSize+len(b))
  iv := ciphertext[:aes.BlockSize]
  if _, err := io.ReadFull(rand.Reader, iv); err != nil {
    return nil, err
  }
  cfb := cipher.NewCFBEncrypter(block, iv)
  cfb.XORKeyStream(ciphertext[aes.BlockSize:], b)
  return ciphertext, nil
}

func Decrypt(key, bytz []byte) ([]byte, error) { 
  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }
  iv := bytz[:aes.BlockSize]
  ciphertext := bytz[aes.BlockSize:]
  cfb := cipher.NewCFBDecrypter(block, iv)
  plaintext := make([]byte, len(ciphertext))
  cfb.XORKeyStream(plaintext, ciphertext)
  data, err := Decode(plaintext)
  if err != nil {
    return nil, err
  }
  return data, nil
}

func Hashkey(key []byte) []byte {
  var cipherkey []byte
  sum := sha256.Sum256(key)
  cipherkey = append(cipherkey, sum[:]...)
  return cipherkey
}

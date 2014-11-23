package fedops_provider

import (
  _"fmt"
  _"net/http"
)

type DigitalOceanAuth struct {
  ApiKey string
}

type DigitalOcean struct {
}

func (d *DigitalOcean) CreateKeypair(Keypair) (string, error) {  
  return "1234-5678-abcd-0000", nil
}

func DigitalOceanProvider(auth DigitalOceanAuth) DigitalOcean {
  return DigitalOcean{}
}
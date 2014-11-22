package fedops_provider

import (
  _"fmt"
)

type DigitalOcean struct {
}

func (d *DigitalOcean) CreateKeypair(Keypair) (error) {
  return nil
}

func DigitalOceanProvider(auth DigitalOceanAuth) DigitalOcean {
  return DigitalOcean{}
}
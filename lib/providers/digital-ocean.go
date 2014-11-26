package fedops_provider

import (
  "fmt"
  "bytes"
  _"io/ioutil"
  "net/http"
  "encoding/json"
  "strconv"
)

type DigitalOceanAuth struct {
  ApiKey string
}

type DigitalOcean struct {
  ApiKey string
  ApiEndpoint string
  KeyURI string
  ImageURI string
  VM_URI string
}

func (digo *DigitalOcean) Name() string {
  return "DigitalOcean"
}

func (digo *DigitalOcean) CreateKeypair(key Keypair) (string, error) {  
  client := &http.Client{}
  //resp, err := client.Get(digo.ApiEndpoint + digo.KeyURI)
  //fmt.Printf("%+v \r\n", key)
  reqJSON := []byte("{\"name\":\"FedOps-ClusterKey-001\", \"public_key\":\"" + string(key.PublicSSH) + " fedops\"}")
  req, err := http.NewRequest("POST", digo.ApiEndpoint + digo.KeyURI, bytes.NewBuffer(reqJSON))
  req.Header.Add("X-FedOps-Provider", digo.Name())
  req.Header.Add("Content-Type", "application/json")
  req.Header.Add("Authorization", "Bearer " + digo.ApiKey)
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  //fmt.Println("Response Status:", resp.Status)
  //fmt.Println("Response Headers:", resp.Header)

  // This is a post so we may have a post Body
  decoder := json.NewDecoder(resp.Body)
  var data interface{}

  err = decoder.Decode(&data)
  if err != nil {
    fmt.Println("JSON body not formated correctly", err.Error())
    return "", err;
  }

  jsonMap := data.(map[string]interface{})
  ssh_key := jsonMap["ssh_key"].(map[string]interface{})

  return strconv.FormatFloat(ssh_key["id"].(float64), 'f', 0, 32), nil
}

func (digo *DigitalOcean) CreateImage() {  
}

func (digo *DigitalOcean) CreateVM() {  
}

func DigitalOceanProvider(auth DigitalOceanAuth) DigitalOcean {
  return DigitalOcean{
    ApiKey: auth.ApiKey,
    ApiEndpoint: "https://api.digitalocean.com",
    KeyURI: "/v2/account/keys",
  }
}
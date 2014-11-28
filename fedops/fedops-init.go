package main

import (
  // Standard
  _"runtime"
  _"os"
  "bufio"
  "fmt"
  "strings"
  // 3rd Party
  "github.com/codegangsta/cli"
  "github.com/gopass"
  // FedOps
  "github.com/FedOps/lib"
)


func commandInit(stdin bufio.Reader, pwd string, cmds []cli.Command) {
  cmd := cli.Command {
    Name: "init",
    ShortName: "i",
    Usage: "create a new cluster",
    Action: func(c *cli.Context) {
      hasConfig := fedops.HasConfigFile(pwd)
      if hasConfig == true {
        fmt.Println("FedOps cluster config file already exist")
        fmt.Println("Try 'fedops info' or 'fedops dashboard'")
        return
      }

      fmt.Println("Providers: digital ocean, aws, google cloud, microsoft azure, openstack")
      fmt.Printf("Enter Cloud Provider... ")
      cloudProvider, _ := stdin.ReadString('\n')
      cloudProvider = strings.Trim(cloudProvider, "\n")

      tokens := fedops.ProviderTokens{}
      switch cloudProvider {
        case "digital ocean":
          fmt.Printf("Enter Digital Ocean API Token... ")
          digoToken, _ := stdin.ReadString('\n')
          digoToken = strings.Trim(digoToken, "\n")
          tokens.AccessToken = digoToken
        case "aws":
          fmt.Printf("Enter Amazon Webservice Token... ")
          awsToken, _ := stdin.ReadString('\n')
          awsToken = strings.Trim(awsToken, "\n")
          tokens.AccessToken = awsToken

          fmt.Printf("Enter Amazon Webservice Security Token... ")
          awsSecurityToken, _ := stdin.ReadString('\n')
          awsSecurityToken = strings.Trim(awsSecurityToken, "\n")
          tokens.SecurityToken = awsSecurityToken
        case "google cloud":
        case "microsoft azure":
        case "openstack":
        default:
          fmt.Println("Unknown provider")
        return
      }

      fmt.Printf("Cluster Config Password... ")
      passwd := gopass.GetPasswd()

      if len(passwd) < MinKeyLength {
        fmt.Printf("Password to short, must be at least %v characters long \r\n", MinKeyLength)
        return
      }

      fmt.Printf("Repeat... ")
      passwdR := gopass.GetPasswd()

      if string(passwd) != string(passwdR) {
        fmt.Println("Passwords don't match")
        return
      }

      if c.Bool("no-harden") == true {
        fmt.Println("WARNING")
        fmt.Println("This cluster's base image will NOT be hardened")
        fmt.Println("Full disk encryption and iptables have been disabled")
      }
      
      fed, err := fedops.CreateDispatcher(passwd, pwd, false)
      if err != nil {
        fmt.Println(err.Error())
        return
      }
      promise := make(chan uint)
      go fed.InitCloudProvider(promise, cloudProvider, tokens)
      status := <- promise
      switch status {
        case fed.Error:
          fmt.Println("Unable to use " + cloudProvider)
        case fed.Ok:
        case fed.Unknown:
          fmt.Println("Unknown")
      }
    },
  }
  cmds = append(cmds, cmd)
}
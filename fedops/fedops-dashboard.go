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


func commandConnect(stdin bufio.Reader, pwd string, cmds []cli.Command) {
  cmd := cli.Command {
    Name: "dashboard",
    ShortName: "d",
    Usage: "create a secure tunnel to the cluster and host a dashboard",
    Action: func(c *cli.Context) {
      hasConfig := fedops.HasConfigFile(pwd)
      if hasConfig == false {
        fmt.Println("FedOps cluster config file does not already exist")
        fmt.Println("Try 'fedops init' or 'fedops connect'")
        return
      }

      fed, err := initFedops(pwd)
      if err != nil {
        fmt.Println("Incorrect Password")
        return
      }
      fmt.Printf("%+v \r\n", fed.Cipherkey)
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      //sessionTasks := []string{""}
      //for _, t := range warehouseTasks {
      //  fmt.Println(t)
      //}
    },
  }
  cmds = append(cmds, cmd)
}
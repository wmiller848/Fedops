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
    Name: "info",
    ShortName: "if",
    Usage: "get info on the cluster",
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

      //fmt.Printf("%+v \r\n", fed.Config)

      //fmt.Println("ClusterID | " + fed.Config.ClusterID)
      fmt.Println("Warehouses")
      if len(fed.Config.Warehouses) > 0 {
        for _, warehouse := range fed.Config.Warehouses {
          fmt.Println("\t", "- ", warehouse.WarehouseID, " ", warehouse.IP, " | ", "uptime...")
        }
      } else {
        fmt.Println("\t", "No warehouses available")
        fmt.Println("\t", "Try 'fedops warehouse create'")
      }

      fmt.Println("Trucks")
      if len(fed.Config.Trucks) > 0 {

      } else {
        fmt.Println("\t", "No trucks available")
        fmt.Println("\t", "Try 'fedops truck create'")
      }
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      //warehouseTasks := []string{"create", "destroy"}
      //for _, t := range warehouseTasks {
      //  fmt.Println(t)
      //}
    },
  }
  cmds = append(cmds, cmd)
}
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
    Name: "warehouse",
    ShortName: "w",
    Usage: "create a new warehouse",
    Action: func(c *cli.Context) {

      args := c.Args()
      if len(args) == 0 {
        fmt.Println("not enought arguments")
        return
      } else if len(args) > 1 {
        fmt.Println("too many arguments")
        return
      }

      provider := "auto"
      memSize := "auto"
      diskSize := "auto"
      numVcpus := "auto"

      cmd := args[0]
      switch cmd {
        case "create":
          fed, err := initFedops(pwd)
          if err != nil {
            fmt.Println("Incorrect Password")
            return
          }
          promise := make(chan uint)
          go fed.CreateTruck(promise, provider, memSize, diskSize, numVcpus)
          status :=  <- promise
          switch status {
            case fed.Error:
              fmt.Println("Error")
            case fed.Ok:
              //fmt.Println("Ok")
            case fed.Unknown:
              fmt.Println("Unknown")
          }
          
        default:
          fmt.Println("Unknown argument for 'fedops warehouse'")
      }
    },
    Flags: []cli.Flag {
      cli.StringFlag {
        Name: "provider",
        Usage: "provider for warehouse, otherwise automatically selects an available provider",
      },
      cli.StringFlag {
        Name: "memory-size",
        Usage: "memory size for warehouse, otherwise automatically selects default for provider",
      },
      cli.StringFlag {
        Name: "disk-size",
        Usage: "disk size for warehouse, otherwise automatically selects default for provider",
      },
      cli.StringFlag {
        Name: "vcpus-size",
        Usage: "number of vcpus for warehouse, otherwise automatically selects default for provider",
      },
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      warehouseTasks := []string{"create", "destroy", "transfer"}
      for _, t := range warehouseTasks {
        fmt.Println(t)
      }
    },
  }
  cmds = append(cmds, cmd)
}
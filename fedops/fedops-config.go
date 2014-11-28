package main

import (
  // Standard
  _"runtime"
  _"os"
  "bufio"
  "fmt"
  _"strings"
  // 3rd Party
  "github.com/codegangsta/cli"
  _"github.com/gopass"
  // FedOps
  _"github.com/FedOps/lib"
)


func commandConfig(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "config",
    ShortName: "cf",
    Usage: "access the cluster config",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Talking to the cloud...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      warehouseTasks := []string{"password", "export"}
      for _, t := range warehouseTasks {
        fmt.Println(t)
      }
    },
  }
  return cmd
}
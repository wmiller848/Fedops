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
    Name: "connect",
    ShortName: "c",
    Usage: "connect to a cluster",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Talking to the cloud...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      warehouseTasks := []string{"create", "destroy"}
      for _, t := range warehouseTasks {
        fmt.Println(t)
      }
    },
  }
  cmds = append(cmds, cmd)
}
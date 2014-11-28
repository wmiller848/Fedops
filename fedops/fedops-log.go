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


func commandLog(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "log",
    ShortName: "l",
    Usage: "output event stream logs",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Doing log thing...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      logTasks := []string{"docker", "system", "network", "warehouse", "truck"}
      for _, t := range logTasks {
        fmt.Println(t)
      }
    },
  }
  return cmd
}
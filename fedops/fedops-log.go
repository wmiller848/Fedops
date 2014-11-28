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
  cmds = append(cmds, cmd)
}
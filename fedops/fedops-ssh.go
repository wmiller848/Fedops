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


func commandSSH(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "ssh",
    ShortName: "s",
    Usage: "ssh into a warehouse or truck",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Doing ssh thing...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      sshTasks := []string{"warehouse", "truck", "keys"}
      for _, t := range sshTasks {
        fmt.Println(t)
      }
    },
  }
  return cmd
}
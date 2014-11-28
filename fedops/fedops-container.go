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


func commandContainer(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "container",
    ShortName: "ct",
    Usage: "create, destroy, assign, ship, or info",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Do the container thing...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      containerTasks := []string{"create", "destroy", "assign", "ship", "info"}
      for _, t := range containerTasks {
        fmt.Println(t)
      }
    },
  }
  return cmd
}
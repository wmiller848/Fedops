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


func commandUse(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "use",
    ShortName: "u",
    Usage: "use a manifest file for the cluster",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Do the use thing...")
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      //for _, t := range tasks {
      //  fmt.Println(t)
      //}
    },
  }
  return cmd
}
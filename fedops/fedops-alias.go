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
  // FedOps
  _"github.com/FedOps/lib"
)


func commandAlias(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "alias",
    ShortName: "a",
    Usage: "alias a truck",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Do the alias thing...")
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
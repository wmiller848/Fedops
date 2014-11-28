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


func commandEnv(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "env",
    ShortName: "e",
    Usage: "create a new deployment env",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Do the env thing...")
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
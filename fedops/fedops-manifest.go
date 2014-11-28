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
    Name: "manifest",
    ShortName: "m",
    Usage: "generate a manifest file from the cluster",
    Action: func(c *cli.Context) {
      //fmt.Printf("%+v \r\n", c)
      fmt.Println("Do the manifest thing...")
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
  cmds = append(cmds, cmd)
}
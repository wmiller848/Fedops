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


func commandManifest(stdin *bufio.Reader, pwd string) cli.Command {
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
  return cmd
}
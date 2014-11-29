// The MIT License (MIT)

// Copyright (c) 2014 William Miller

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
  "github.com/FedOps/lib"
)


func commandDashboard(stdin *bufio.Reader, pwd string) cli.Command {
  cmd := cli.Command {
    Name: "dashboard",
    ShortName: "d",
    Usage: "create a secure tunnel to the cluster and host a dashboard",
    Action: func(c *cli.Context) {
      hasConfig := fedops.HasConfigFile(pwd)
      if hasConfig == false {
        fmt.Println("FedOps cluster config file does not already exist")
        fmt.Println("Try 'fedops init' or 'fedops connect'")
        return
      }

      fed, err := initFedops(pwd)
      if err != nil {
        fmt.Println("Incorrect Password")
        return
      }
      fmt.Printf("%+v \r\n", fed.Cipherkey)
    },
    BashComplete: func(c *cli.Context) {
      // This will complete if no args are passed
      if len(c.Args()) > 0 {
        return
      }
      //sessionTasks := []string{""}
      //for _, t := range warehouseTasks {
      //  fmt.Println(t)
      //}
    },
  }
  return cmd
}
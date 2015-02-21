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
	"bufio"
	"fmt"
	// 3rd Party
	"github.com/codegangsta/cli"
	// FedOps
	"github.com/FedOps/lib"
)

func commandWarehouse(stdin *bufio.Reader, pwd string) cli.Command {
	cmd := cli.Command{
		Name:      "warehouse",
		// ShortName: "w",
		Usage:     "manage fedops warehouses: create, destroy",
    Subcommands: []cli.Command{
      cli.Command{
        Name: "create",
        Action: func(c *cli.Context) {
          providerName := "auto"
          memSize := "auto"
          diskSize := "auto"
          numVcpus := "auto"

          fed, err := initFedops(pwd)
          if err != nil {
            fmt.Println("Incorrect Password")
            return
          }
          promise := make(chan fedops.FedopsAction)
          go fed.CreateWarehouse(promise, providerName, memSize, diskSize, numVcpus)
          result := <- promise
          switch result.Status {
          case fedops.FedopsError:
            fmt.Println("Error")
          case fedops.FedopsOk:
            //fmt.Println("Ok")
          case fedops.FedopsUnknown:
            fmt.Println("Unknown")
          }
        },
      },
      cli.Command{
        Name: "destroy",
        Action: func(c *cli.Context) {
          fed, err := initFedops(pwd)
          if err != nil {
            fmt.Println("Incorrect Password")
            return
          }

          if len (c.Args()) == 0 {
            fmt.Println("Supply a warehouse ID") 
            return
          }

          warehouseID := c.Args()[0] //c.String("warehouseID")
          if warehouseID == "" {
            fmt.Println("Supply a warehouse ID")
            return
          }

          promise := make(chan fedops.FedopsAction)
          go fed.DestroyWarehouse(promise, warehouseID)
          result := <- promise
          switch result.Status {
          case fedops.FedopsError:
            // fmt.Println("Error")
          case fedops.FedopsOk:
            //fmt.Println("Ok")
          case fedops.FedopsUnknown:
            fmt.Println("Unknown")
          }
        },
      },
    },
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "provider",
				Usage: "create : provider for warehouse, otherwise automatically selects an available provider",
			},
			cli.StringFlag{
				Name:  "memory-size",
				Usage: "create : memory size for warehouse, otherwise automatically selects default for provider",
			},
			cli.StringFlag{
				Name:  "disk-size",
				Usage: "create : disk size for warehouse, otherwise automatically selects default for provider",
			},
			cli.StringFlag{
				Name:  "vcpus-size",
				Usage: "create : number of vcpus for warehouse, otherwise automatically selects default for provider",
			},
      cli.StringFlag{
        Name:  "warehouseID",
        Usage: "destory : warehouse ID",
      },
		},
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if len(c.Args()) > 0 {
				return
			}
			warehouseTasks := []string{"create", "destroy", "transfer"}
			for _, t := range warehouseTasks {
				fmt.Println(t)
			}
		},
	}
	return cmd
}

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
	_ "os"
	_ "runtime"
	_ "strings"
	// 3rd Party
	"github.com/codegangsta/cli"
	_ "github.com/gopass"
	// FedOps
	"github.com/FedOps/lib"
)

func commandInfo(stdin *bufio.Reader, pwd string) cli.Command {
	cmd := cli.Command{
		Name:      "info",
		ShortName: "if",
		Usage:     "get info on the cluster",
		Action: func(c *cli.Context) {
			hasConfig := fedops.HasConfigFile(pwd)
			if hasConfig == false {
				fmt.Println("Fedops cluster config file does not already exist")
				fmt.Println("Try 'fedops init' or 'fedops connect'")
				return
			}

			fed, err := initFedops(pwd)
			if err != nil {
				fmt.Println("Incorrect Password")
				return
			}

			promise := make(chan fedops.FedopsAction)
			go fed.Refresh(promise)
			result := <-promise
			switch result.Status {
			case fedops.FedopsError:
				fmt.Println("Error")
			case fedops.FedopsOk:
				//fmt.Println("Ok")
			case fedops.FedopsUnknown:
				fmt.Println("Unknown")
			}
			//fmt.Printf("%+v \r\n", fed.Config)

			//fmt.Println("ClusterID | " + fed.Config.ClusterID)
			fmt.Println("Warehouses")
			if len(fed.Config.Warehouses) > 0 {
				for _, warehouse := range fed.Config.Warehouses {
					fmt.Println("\t", " -", warehouse.WarehouseID, " - ", warehouse.IPV4, " | ", warehouse.Status)
				}
			} else {
				fmt.Println("\t", "No warehouses available")
				fmt.Println("\t", "Try 'fedops warehouse create'")
			}

			fmt.Println("Trucks")
			if len(fed.Config.Trucks) > 0 {
				for _, truck := range fed.Config.Trucks {
					fmt.Println("\t", " -", truck.TruckID, " - ", truck.IPV4, " | ", truck.Status)
				}
			} else {
				fmt.Println("\t", "No trucks available")
				fmt.Println("\t", "Try 'fedops truck create'")
			}
		},
		BashComplete: func(c *cli.Context) {
			// This will complete if no args are passed
			if len(c.Args()) > 0 {
				return
			}
			//warehouseTasks := []string{"create", "destroy"}
			//for _, t := range warehouseTasks {
			//  fmt.Println(t)
			//}
		},
	}
	return cmd
}

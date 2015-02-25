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
	"strings"
	// 3rd Party
	"github.com/codegangsta/cli"
	"code.google.com/p/gopass"
	// FedOps
	"github.com/Fedops/lib"
)

func commandInit(stdin *bufio.Reader, pwd string) cli.Command {
	cmd := cli.Command{
		Name:      "init",
		// ShortName: "i",
		Usage:     "create a new cluster",
		Action: func(c *cli.Context) {
			hasConfig := fedops.HasConfigFile(pwd)
			if hasConfig == true {
				fmt.Println("FedOps cluster config file already exist")
				fmt.Println("Try 'fedops info' or 'fedops dashboard'")
				return
			}

			fmt.Println("Providers: digital ocean, aws, google cloud, microsoft azure, openstack")
			fmt.Printf("Enter Cloud Provider... ")
			cloudProvider, _ := stdin.ReadString('\n')
			cloudProvider = strings.Trim(cloudProvider, "\n")

			tokens := fedops.Tokens{}
			switch cloudProvider {
			case "digital ocean":
				fmt.Printf("Enter Digital Ocean API Token... ")
				digoToken, _ := stdin.ReadString('\n')
				digoToken = strings.Trim(digoToken, "\n")
				tokens.AccessToken = digoToken
			case "aws":
				fmt.Printf("Enter Amazon Webservice Token... ")
				awsToken, _ := stdin.ReadString('\n')
				awsToken = strings.Trim(awsToken, "\n")
				tokens.AccessToken = awsToken

				fmt.Printf("Enter Amazon Webservice Security Token... ")
				awsSecurityToken, _ := stdin.ReadString('\n')
				awsSecurityToken = strings.Trim(awsSecurityToken, "\n")
				tokens.SecurityToken = awsSecurityToken
			case "google cloud":
			case "microsoft azure":
			case "openstack":
			default:
				fmt.Println("Unknown provider")
				return
			}

			fmt.Printf("Cluster Config Password... ")
			passwd, _ := gopass.GetPass("")

			if len(passwd) < MinKeyLength {
				fmt.Printf("Password to short, must be at least %v characters long \r\n", MinKeyLength)
				return
			}

			fmt.Printf("Repeat... ")
			passwdR, _ := gopass.GetPass("")

			if string(passwd) != string(passwdR) {
				fmt.Println("Passwords don't match")
				return
			}

			if c.Bool("no-harden") == true {
				fmt.Println("WARNING")
				fmt.Println("This cluster's base image will NOT be hardened")
				fmt.Println("Full disk encryption and iptables have been disabled")
			}

			fed, err := fedops.CreateDispatcher([]byte(passwd), pwd, false)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			promise := make(chan fedops.FedopsAction)
			go fed.InitCloudProvider(promise, cloudProvider, tokens)
			result := <-promise
			switch result.Status {
			case fedops.FedopsError:
				fmt.Println("Unable to use " + cloudProvider)
			case fedops.FedopsOk:
			case fedops.FedopsUnknown:
				fmt.Println("Unknown")
			}
		},
	}
	return cmd
}

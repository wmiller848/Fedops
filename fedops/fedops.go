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
	"os"
	"runtime"
	// 3rd Party
	"code.google.com/p/gopass"
	"github.com/codegangsta/cli"
	// FedOps
	"github.com/wmiller848/Fedops/lib/dispatcher"
	"github.com/wmiller848/Fedops/lib/encryption"
)

const (
	MinKeyLength int = 6
)

func initFedops(pwd string) (*fedops.Dispatcher, error) {
	var passwd []byte
	session_key := os.Getenv("FEDOPS_SESSION_KEY")
	useSession := false
	if session_key != "" {
		// The user has set the session key
		// which is just the encoded cipherkey
		var err error
		passwd, err = fedops_encryption.Decode([]byte(session_key))
		if err != nil {
			return nil, err
		}
		useSession = true
	} else {
		fmt.Printf("Cluster Config Password... ")
		pass, _ := gopass.GetPass("")
		passwd = []byte(pass)
	}

	fed, err := fedops.CreateDispatcher(passwd, pwd, useSession)
	if err != nil {
		return nil, err
	}
	return fed, nil
}

//
//
func main() {
	//
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)

	pwd := os.Getenv("PWD")
	stdin := bufio.NewReader(os.Stdin)

	_cli := cli.NewApp()
	_cli.Name = "FedOps"
	_cli.Author = "W. Chase Miller"
	_cli.Email = "wmiller.fedops@gmail.com"
	_cli.Usage = "Docker continuous deployment and cloud management, see https://github.com/wmiller848/Fedops for guides"
	_cli.Version = "0.0.1"
	_cli.EnableBashCompletion = true

	commands := []cli.Command{}

	// Register subcommands
	commands = append(commands, commandConfig(stdin, pwd))
	commands = append(commands, commandContainer(stdin, pwd))
	commands = append(commands, commandInfo(stdin, pwd))
	commands = append(commands, commandInit(stdin, pwd))
	commands = append(commands, commandLog(stdin, pwd))
	commands = append(commands, commandSession(stdin, pwd))
	commands = append(commands, commandSSH(stdin, pwd))
	commands = append(commands, commandTruck(stdin, pwd))
	commands = append(commands, commandWarehouse(stdin, pwd))

	_cli.Commands = commands

	_cli.Action = func(c *cli.Context) {
		fmt.Println("Fedops is a tool for continuous deployment")
		fmt.Println("Try 'fedops help'")
	}
	_cli.Run(os.Args)
}

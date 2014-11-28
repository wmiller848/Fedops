package main

import (
	// Standard
  "runtime"
  "os"
  "bufio"
	"fmt"
  _"strings"
	// 3rd Party
	"github.com/codegangsta/cli"
  "github.com/gopass"
	// FedOps
  "github.com/FedOps/lib"
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
    passwd, err = fedops.Decode([]byte(session_key))
    if err != nil {
      return nil, err
    }
    useSession = true
  } else {
    fmt.Printf("Cluster Config Password... ")
    passwd = gopass.GetPasswd()
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
	_cli.Usage = "Docker continuous deployment made easy, see https://github.com/wmiller848/FedOps for guides"
  _cli.Version = "0.0.1"
  _cli.EnableBashCompletion = true

  commands := []cli.Command{}

  commands = append(commands, commandAlias(stdin, pwd))
  commands = append(commands, commandConfig(stdin, pwd))
  commands = append(commands, commandConnect(stdin, pwd))
  commands = append(commands, commandContainer(stdin, pwd))
  commands = append(commands, commandDashboard(stdin, pwd))
  commands = append(commands, commandEnv(stdin, pwd))
  commands = append(commands, commandInfo(stdin, pwd))
  commands = append(commands, commandInit(stdin, pwd))
  commands = append(commands, commandLog(stdin, pwd))
  commands = append(commands, commandManifest(stdin, pwd))
  commands = append(commands, commandSession(stdin, pwd))
  commands = append(commands, commandSSH(stdin, pwd))
  commands = append(commands, commandTruck(stdin, pwd))
  commands = append(commands, commandUse(stdin, pwd))
  commands = append(commands, commandWarehouse(stdin, pwd))

  _cli.Commands = commands

	_cli.Action = func(c *cli.Context) {
		fmt.Println("Fedops is a tool for continuous deployment")
    fmt.Println("Try 'fedops help'")
	}
	_cli.Run(os.Args)
}

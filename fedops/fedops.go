package main

import (
	// Standard
  "runtime"
  "os"
  "bufio"
	"fmt"
  "strings"
	// 3rd Party
	"github.com/codegangsta/cli"
  "github.com/gopass"
	// FedOps
  "github.com/FedOps/lib"
)

//
//
func main() {
  //
  numCpus := runtime.NumCPU()
  runtime.GOMAXPROCS(numCpus)

  env := os.Environ()
  pwd := ""
  for _, e := range env {
    pair := strings.Split(e, "=")
    name := pair[0]
    value := pair[1]
    if name == "PWD" {
      pwd = value
    }
  }
  stdin := bufio.NewReader(os.Stdin)
	_cli := cli.NewApp()
	_cli.Name = "FedOps"
	_cli.Usage = "fedops init"
  _cli.EnableBashCompletion = true
	_cli.Commands = []cli.Command {
    {
      Name: "init",
    	ShortName: "i",
    	Usage: "create a new cluster",
    	Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Providers: digital ocean, aws, google cloud, microsoft azure, openstack")
        fmt.Printf("Enter Cloud Provider... ")
        cloudProvider, _ := stdin.ReadString('\n')
        cloudProvider = strings.Trim(cloudProvider, "\n")

        tokens := fedops.ProviderTokens{}
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

        fmt.Printf("Cluster Password... ")
        passwd := string(gopass.GetPasswd())
        fmt.Printf("Repeat... ")
        passwdR := string(gopass.GetPasswd())

        if passwd != passwdR {
          fmt.Println("Passwords don't match")
          return
        }
        
        fedops, loaded := fedops.CreateDispatcher(passwd, pwd)
        if loaded == true {
          fmt.Println("FedOps Cluster Config file found")
          return
        }
        promise := make(chan uint)
        var status uint
        fedops.InitCloudProvider(promise, cloudProvider, tokens)
        status = <- promise
        switch status {
          case fedops.Error:
            fmt.Println("Unable to use " + cloudProvider)
          case fedops.Ok:
            //fmt.Println("")
          case fedops.Unknown:
            fmt.Println("Unknown")
        }

        //fmt.Println("Configuring a new fedops cluster...")

        //fmt.Println("Running Test Transaction...")

        //fmt.Println("Creating a new warehouse...")
        //fmt.Println("Shipping to warehouse...")

        //fmt.Println("Creating a new truck...")
        //fmt.Println("Shipping to truck...")
    	},
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        initTasks := []string{"--manifest", "--no-harden"}
        for _, t := range initTasks {
          fmt.Println(t)
        }
      },
  	},
    {
      Name: "connect",
      ShortName: "con",
      Usage: "connect to a cluster",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        warehouseTasks := []string{"create", "destroy"}
        for _, t := range warehouseTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "info",
      ShortName: "i",
      Usage: "get info on the cluster",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
    {
      Name: "session",
      ShortName: "session",
      Usage: "create a new session, required for every use",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Decrypt Fedops Config")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        warehouseTasks := []string{"create", "destroy"}
        for _, t := range warehouseTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "warehouse",
      ShortName: "w",
      Usage: "create a new warehouse machine",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        warehouseTasks := []string{"create", "destroy"}
        for _, t := range warehouseTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "truck",
      ShortName: "t",
      Usage: "create a new truck machine",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
    {
      Name: "container",
      ShortName: "c",
      Usage: "create, destroy, assign, ship, or list",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        containerTasks := []string{"create", "destroy", "assign", "ship", "list"}
        for _, t := range containerTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "env",
      ShortName: "e",
      Usage: "create a new deployment env",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
    {
      Name: "fork",
      ShortName: "f",
      Usage: "fork a truck machine",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
    {
      Name: "use",
      ShortName: "u",
      Usage: "use a manifest file for the cluster",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
    {
      Name: "manifest",
      ShortName: "m",
      Usage: "generate a manifest file for the cluster",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
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
    },
	}

	_cli.Action = func(c *cli.Context) {
		fmt.Println("Please enter a command...")
	}
	_cli.Run(os.Args)
}

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
  var pwd string
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
        hasConfig := fedops.HasConfigFile(pwd)
        if hasConfig == true {
          fmt.Println("FedOps cluster config file already exist")
          fmt.Println("Try 'fedops info' or 'fedops dashboard'")
          return
        }
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
        passwd := gopass.GetPasswd()

        if len(passwd) < 4 {
          fmt.Println("Password to short, must be at least 4 characters long")
          return
        }

        fmt.Printf("Repeat... ")
        passwdR := gopass.GetPasswd()

        if string(passwd) != string(passwdR) {
          fmt.Println("Passwords don't match")
          return
        }

        if c.Bool("no-harden") == true {
          fmt.Println("You have requested image full disk encryption be disabled!")
        }
        
        fed, _ := fedops.CreateDispatcher(passwd, pwd, false)
        promise := make(chan uint)
        var status uint
        fed.InitCloudProvider(promise, cloudProvider, tokens)
        status = <- promise
        switch status {
          case fed.Error:
            fmt.Println("Unable to use " + cloudProvider)
          case fed.Ok:
            //fmt.Println("")
          case fed.Unknown:
            fmt.Println("Unknown")
        }
    	},
      Flags: []cli.Flag {
        cli.BoolFlag {
          Name: "no-harden",
          Usage: "enable or disable full disk encryption",
        },
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        initTasks := []string{"manifest", "no-harden"}
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
        hasConfig := fedops.HasConfigFile(pwd)
        if hasConfig == false {
          fmt.Println("FedOps cluster config file does not already exist")
          fmt.Println("Try 'fedops init' or 'fedops connect'")
          return
        }

        var passwd []byte
        session_key := os.Getenv("FEDOPS_SESSION_KEY")
        useSession := false
        if session_key != "" {
          var err error
          passwd, err = fedops.Decode([]byte(session_key))
          if err != nil {
            fmt.Println(err.Error())
            return
          }
          useSession = true
        } else {
          fmt.Printf("Cluster Password... ")
          passwd = gopass.GetPasswd()
        }

        fed, err := fedops.CreateDispatcher(passwd, pwd, useSession)
        if err != nil {
          fmt.Println("Incorrect Password")
          return
        }
        fmt.Println("ClusterID | " + fed.Config.ClusterID)
        fmt.Println("warehouses")
        fmt.Println("\tn/a")
        fmt.Println("trucks")
        fmt.Println("\tn/a")
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
      Usage: "output session key",
      Action: func(c *cli.Context) {
        hasConfig := fedops.HasConfigFile(pwd)
        if hasConfig == false {
          fmt.Println("FedOps cluster config file does not already exist")
          fmt.Println("Try 'fedops init' or 'fedops connect'")
          return
        }

        fmt.Printf("Cluster Password... ")
        passwd := gopass.GetPasswd()
        fed, err := fedops.CreateDispatcher(passwd, pwd, false)
        if err != nil {
          fmt.Println("Incorrect Password")
          return
        }
        fmt.Println("export FEDOPS_SESSION_KEY=" + string(fed.Key))
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        warehouseTasks := []string{"create"}
        for _, t := range warehouseTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "warehouse",
      ShortName: "w",
      Usage: "create a new warehouse",
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
      Usage: "create a new truck",
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
      Name: "alias",
      ShortName: "a",
      Usage: "alias a truck",
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
      Usage: "use a manifest file for the cluster [WARNING]",
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
      Usage: "generate a manifest file from the cluster",
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

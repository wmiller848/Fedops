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

  commands := []cli.Command
  commandInit(stdin, pwd, commands)
  commandConnect(stdin, pwd, commands)

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

        fmt.Printf("Cluster Config Password... ")
        passwd := gopass.GetPasswd()

        if len(passwd) < MinKeyLength {
          fmt.Printf("Password to short, must be at least %v characters long \r\n", MinKeyLength)
          return
        }

        fmt.Printf("Repeat... ")
        passwdR := gopass.GetPasswd()

        if string(passwd) != string(passwdR) {
          fmt.Println("Passwords don't match")
          return
        }

        if c.Bool("no-harden") == true {
          fmt.Println("WARNING")
          fmt.Println("This cluster's base image will NOT be hardened")
          fmt.Println("Full disk encryption and iptables have been disabled")
        }
        
        fed, err := fedops.CreateDispatcher(passwd, pwd, false)
        if err != nil {
          fmt.Println(err.Error())
          return
        }
        promise := make(chan uint)
        go fed.InitCloudProvider(promise, cloudProvider, tokens)
        status := <- promise
        switch status {
          case fed.Error:
            fmt.Println("Unable to use " + cloudProvider)
          case fed.Ok:
          case fed.Unknown:
            fmt.Println("Unknown")
        }
    	},
      Flags: []cli.Flag {
        cli.BoolFlag {
          Name: "no-events",
          Usage: "disable event stream, which via 'fedops log' provides access to docker, system, and network information",
        },
        cli.BoolFlag {
          Name: "no-harden",
          Usage: "disable full disk encryption and iptables",
        },
        cli.StringFlag {
          Name: "bootstrap",
          Usage: "path or url to shell script to execute at the end of the defualt setup",
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
      ShortName: "c",
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
      Name: "config",
      ShortName: "cf",
      Usage: "access the cluster config",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Talking to the cloud...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        warehouseTasks := []string{"password", "export"}
        for _, t := range warehouseTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "info",
      ShortName: "if",
      Usage: "get info on the cluster",
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

        //fmt.Printf("%+v \r\n", fed.Config)

        //fmt.Println("ClusterID | " + fed.Config.ClusterID)
        fmt.Println("Warehouses")
        if len(fed.Config.Warehouses) > 0 {
          for _, warehouse := range fed.Config.Warehouses {
            fmt.Println("\t", "- ", warehouse.WarehouseID, " ", warehouse.IP, " | ", "uptime...")
          }
        } else {
          fmt.Println("\t", "No warehouses available")
          fmt.Println("\t", "Try 'fedops warehouse create'")
        }

        fmt.Println("Trucks")
        if len(fed.Config.Trucks) > 0 {

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
    },
    {
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
    },
    {
      Name: "ssh",
      ShortName: "s",
      Usage: "ssh into a warehouse or truck",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Doing ssh thing...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        sshTasks := []string{"warehouse", "truck", "keys"}
        for _, t := range sshTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "log",
      ShortName: "l",
      Usage: "output event stream logs",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Doing log thing...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        logTasks := []string{"docker", "system", "network", "warehouse", "truck"}
        for _, t := range logTasks {
          fmt.Println(t)
        }
      },
    },
    {
      Name: "session",
      ShortName: "se",
      Usage: "output session key",
      Action: func(c *cli.Context) {
        hasConfig := fedops.HasConfigFile(pwd)
        if hasConfig == false {
          fmt.Println("FedOps cluster config file does not already exist")
          fmt.Println("Try 'fedops init' or 'fedops connect'")
          return
        }

        session_key := os.Getenv("FEDOPS_SESSION_KEY")
        if session_key == "" {
          fmt.Println("WARNING")
          fmt.Println("This command will export the encoded key that encypts this config, use good judgment when using a session")
          fmt.Printf("Cluster Config Password... ")
          passwd := gopass.GetPasswd()
          fed, err := fedops.CreateDispatcher(passwd, pwd, false)
          if err != nil {
            fmt.Println("Incorrect Password")
            return
          }
          fmt.Println("export FEDOPS_SESSION_KEY=" + string(fed.Cipherkey))
        } else {
          fmt.Println("FEDOPS_SESSION_KEY is already set to " + session_key)
          fmt.Println("To unset run")
          fmt.Println("export FEDOPS_SESSION_KEY=")
        }
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        //sessionTasks := []string{""}
        //for _, t := range sessionTasks {
        //  fmt.Println(t)
        //}
      },
    },
    {
      Name: "warehouse",
      ShortName: "w",
      Usage: "create a new warehouse",
      Action: func(c *cli.Context) {

        args := c.Args()
        if len(args) == 0 {
          fmt.Println("not enought arguments")
          return
        } else if len(args) > 1 {
          fmt.Println("too many arguments")
          return
        }

        provider := "auto"
        memSize := "auto"
        diskSize := "auto"
        numVcpus := "auto"

        cmd := args[0]
        switch cmd {
          case "create":
            fed, err := initFedops(pwd)
            if err != nil {
              fmt.Println("Incorrect Password")
              return
            }
            promise := make(chan uint)
            go fed.CreateTruck(promise, provider, memSize, diskSize, numVcpus)
            status :=  <- promise
            switch status {
              case fed.Error:
                fmt.Println("Error")
              case fed.Ok:
                //fmt.Println("Ok")
              case fed.Unknown:
                fmt.Println("Unknown")
            }
            
          default:
            fmt.Println("Unknown argument for 'fedops warehouse'")
        }
      },
      Flags: []cli.Flag {
        cli.StringFlag {
          Name: "provider",
          Usage: "provider for warehouse, otherwise automatically selects an available provider",
        },
        cli.StringFlag {
          Name: "memory-size",
          Usage: "memory size for warehouse, otherwise automatically selects default for provider",
        },
        cli.StringFlag {
          Name: "disk-size",
          Usage: "disk size for warehouse, otherwise automatically selects default for provider",
        },
        cli.StringFlag {
          Name: "vcpus-size",
          Usage: "number of vcpus for warehouse, otherwise automatically selects default for provider",
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
    },
    {
      Name: "truck",
      ShortName: "t",
      Usage: "create a new truck",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Do the truck thing...")
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
      ShortName: "ct",
      Usage: "create, destroy, assign, ship, or info",
      Action: func(c *cli.Context) {
        //fmt.Printf("%+v \r\n", c)
        fmt.Println("Do the container thing...")
      },
      BashComplete: func(c *cli.Context) {
        // This will complete if no args are passed
        if len(c.Args()) > 0 {
          return
        }
        containerTasks := []string{"create", "destroy", "assign", "ship", "info"}
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
        fmt.Println("Do the env thing...")
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
        fmt.Println("Do the alias thing...")
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
        fmt.Println("Do the use thing...")
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
    },
	}

	_cli.Action = func(c *cli.Context) {
		fmt.Println("Fedops is a tool for continuous deployment")
    fmt.Println("Try 'fedops help'")
	}
	_cli.Run(os.Args)
}

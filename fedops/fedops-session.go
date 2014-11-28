package main

import (
  // Standard
  _"runtime"
  _"os"
  "bufio"
  "fmt"
  "strings"
  // 3rd Party
  "github.com/codegangsta/cli"
  "github.com/gopass"
  // FedOps
  "github.com/FedOps/lib"
)


func commandConnect(stdin bufio.Reader, pwd string, cmds []cli.Command) {
  cmd := cli.Command {
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
  }
  cmds = append(cmds, cmd)
}
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
	"fmt"
	"runtime"
	//
	"github.com/wmiller848/Fedops/lib/engine/warehouse"
)

func main() {
	numCpus := runtime.NumCPU()
	runtime.GOMAXPROCS(numCpus)

	daemon := fedops_warehouse.CreateDaemon()
	statusChan := make(chan error)
	if daemon != nil {
		go daemon.Listen(statusChan)
		go daemon.StartEventEngine(statusChan)
	}
	err := <-statusChan
	fmt.Println(err.Error())
	// server cert is self signed -> server_cert == ca_cert
	// CA_Pool := x509.NewCertPool()
	// severCert, err := ioutil.ReadFile("./cert.pem")
	// if err != nil {
	//     log.Fatal("Could not load server certificate!")
	// }
	// CA_Pool.AppendCertsFromPEM(severCert)

	// config := tls.Config{RootCAs: CA_Pool}

	// conn, err := tls.Dial("tcp", "127.0.0.1:1337", &config)
	// if err != nil {
	//     log.Fatalf("client: dial: %s", err)
	// }

	// pwd := os.Getenv("PWD")
	// hasConfig := fedops.HasConfigFile(pwd)
	// if hasConfig == false {
	// 	fmt.Println("FedOps cluster config file does not exist")
	// 	return
	// }
	//
	// session_key := os.Getenv("FEDOPS_SESSION_KEY")
	// if session_key == "" {
	// 	fmt.Println("'FEDOPS_SESSION_KEY' enviroment variable has been unset, please reset it to match the session key for your cluster")
	// 	return
	// }
	//
	// key, err := fedops_encryption.Decode([]byte(session_key))
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// fed, err := fedops.CreateDispatcher(key, pwd, true)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// fed_certs := fed.Config.Certs
	//
	// // cert, err := tls.LoadX509KeyPair("./cert.pem", "./key.pem")
	// cert, err := tls.X509KeyPair(fed_certs[0].CertificatePem, fed_certs[0].PrivatePem)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// config := tls.Config{Certificates: []tls.Certificate{cert}}
	// listener, err := tls.Listen("tcp", ":13371", &config)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	//
	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		fmt.Println("server: accept: %s", err)
	// 		break
	// 	}
	// 	fmt.Println("server: accepted from %s", conn.RemoteAddr())
	// 	// go handleConnection(conn)
	// }
}

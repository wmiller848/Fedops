# Fedops #

### IN FREQUENT DEVELOPMENT ###
### This project is not ready for usage :p ###

### Development, Deployment, Done ###

## About ##

Fedops is a cloud vps manager and continuous integration tool designed to easily and securely manage shipping docker containers. Fedops integrates very tightly with docker, and the docker registry.

You can use Fedops to create a cluster of machines that will poll source control, build, test, and deploy to any number of machines you define.

The developer is referred to as the dispatcher

Machines that manage continuous integration are called warehouses

Machines that receive and run shipped containers are called trucks

Fedops believes in robust security, the 'init' command will walk you through configuring a new cluster.
After you have configured the new cluster Fedops will create an encrypted configuration file for the cluster.
This is important as your configuration file will store sensitive information about your account for a given cloud provider, ssh key pairs, and tls certificate key pairs.

You should treat the 'Fedops' file with care.

## Usage ##
Fedops operates like git, vagrant, or grunt in the sense that it is directory based, see usage for an example

### Setup ###
Move into a clean folder
```
cd ~/clusters/example
```

### Create a new cluster ###
```
fedops init
```
Follow the prompts
```
ls -la ./
```
You should see two files
```
Fedops
.fedops-salt
```

### Create a new warehouse ###
```
fedops warehouse create
```

### Create a new truck ###
```
fedops truck create
```

### Create a new container ###
```
fedops container create https://github.com/wmiller848/amazing_example
```

### Ship a container for continuous deployment ###
```
fedops warehouse ship [warehouseID] [containerID]
```

### Deliver a container for running ###
```
fedops truck deliver [truckID] [containerID]
```

### List cluster information ###
```
fedops info
```

The output would look something like this
```
Warehouses
  - 1jbn891h81h01ndh81h - 201.0.10.10 | Up | 3 minutes
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | X, Building... 1 minute ago

Trucks
  - 891h91h981h809hd819 - 201.0.10.11 | up | 1 minute
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | X, Waiting... 1 minute ago

Unshipped Containers
  none
```

If we looked after the build finished we would see something like this
```
Warehouses
  - 1jbn891h81h01ndh81h - 201.0.10.10 | Up | 15 minutes
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | √, Built... 5 minutes ago

Trucks
  - 891h91h981h809hd819 - 201.0.10.11 | up | 10 minutes
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | √, Running... 5 minutes ago

Unshipped Containers
  none
```

Ship it to another truck
```
fedops truck create
fedops truck deliver [truckID] [containerID]
```

If we looked now
```
Warehouses
  - 1jbn891h81h01ndh81h - 201.0.10.10 | Up | 15 minutes
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | √, Built... 10 minutes ago

Trucks
  - 891h91h981h809hd819 - 201.0.10.11 | up | 10 minutes
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | √, Running... 10 minutes ago
  - jsysa819bdoi18h0hd0 - 201.0.10.15 | up | 1 minute
    * 98cn1oh901h109h19h0 - https://github.com/wmiller848/amazing_example | √, Running... 1 minute ago

Unshipped Containers
  none
```

### SSH into a machine ###
```
fedops ssh [warehouseID or truckID]
```

### Help ###
```
fedops help
fedops help [subcommand]
```

## Architecture ##

Fedops uses a push based architecture, after a warehouse is established it will poll source control. From there it will push events to the relevant trucks via secure tcp listeners.

By default all nodes on any provider are fedora images with SELinux installed/enabled

See https://getfedora.org/

See http://en.wikipedia.org/wiki/Security-Enhanced_Linux

Bootstrapping is done at machine creation time and uses ssh for file transfer and passing commands.
After a fedops node comes online communication is secured via a tls tcp connection, the clusterID is used to further encrypt all commands as an additonal layer on top of tls.

Fedops makes use of golang's standard implementations of crypto and ssh

See https://godoc.org/golang.org/x/crypto

Fedops produces '521 Elliptic Curve' based private keys that are equilvant to 15360 bit RSA keys

See https://www.nsa.gov/business/programs/elliptic_curve.shtml

ALL DATA at rest and/or on the wire is encypted at least once in additon to any transport encryption

Fedops does not protect you from running services with known exploits, not updating (although this may be supported in the future) your servers software, metadata collection from large nation-state organizations and/or internet service providers, rootkits or malware that may infect your vps's. Fedops makes it easy to deploy applications and servers, we recommend you frequently cycle your servers (i.e. destroy and create new ones).

## Liceneces ##

The MIT License (MIT)

Copyright (c) William Miller
Last revised 2015

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

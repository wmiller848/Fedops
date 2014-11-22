# FedOps #

### Development, Deployment, Done ###

Fedops is a cloud vps manager and contionous intigration tool designed to easilty and securely manage shipping docker containers. FedOps initigrates very tightly with docker, and the docker registery.

You can use FedOps to create a cluster of machines that will poll source control, build, test, and deploy to any number of enviroments you define.

The developer is referred to as the dispatcher

Machines that manage continous intigration are called warehouses

Machines that receive and run shipped containers are called trucks

FedOps believes in robust security, the 'init' command will walk you through configuring a new cluster.
After you have configured the new cluster FedOps will create an encrypted configuration file for the cluster.
This is important as your configuration file will store sensative information about your account for a given cloud provider.
You should treat the .fedops folder with care, this contains vital information on your stack and should not be deleted manualy.

FedOps operates like git, vagrant, or grunt in the sense that it is directory based, see usage for an example

FedOps uses 3 types of configurations files

1) The .fedops file is your encrpted cluster file, you should never manualy modify this file
  It contains the static ip addresses for all machines in the cluster, the api access keys for a given provider, the ssh keys for access to the cluster

2) A 'some.service' file should contain a json object describing the docker containers that make up your service
  The service file can include a host system file path to the docker files or a http endpoint to retreive the docker file.
  For example you have a key/value storage api, your service file would specify a docker container for perhapse the http service
  and a dockerfile for the database service backing it. Together those two components describe your service.

3) The 'fedops.manifest' file is the json object that describes how your services should be shipped via avaliable trucks
  For example you have two services, a key/value api and a static website. You have defined these two services as two seperate service files, key_value.service & website.service. In the manifest file we specify that we want 2 instances of the website.service and 2 instances of the key_value.service shipped. Fedops will manage applying theses changes to the cluster.

## Usage ##

### From Existing Cluster ###

* cd ~/clusters/wordpress
* fedops connect [fedops_warehouse_ip]
  # Verify ~/.ssh/id_rsa.pub
  # Enter authorization key
  > auth_key
  # enter encryption key
  > fedops_key
  # enter digital ocean api key
  > api_key


// Cluster with one envirment
* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.10 | Up 3 months & 7 days
  #    * FedOps[x01] : "https://githib.com/wmiller/FedOps" | √, Built...1 minute ago
  #    * AmazingOtherApp[x02] : "https://githib.com/AmazingUser/AmazingOtherApp" | √, Built...2 hours ago
  #      qzj1h8o -> trucks.www1, trucks.www2, trucks.www3
  #  - cd-uswest2 - 201.0.10.11 | Up 3 months & 7 days
  #    * FedOps-www[x01] : "https://githib.com/wmiller/FedOps-www" | √, Built...2 hours ago
  #      s819af1 -> trucks.www1, trucks.www2
  #
  # trucks
  #  - www1 - 201.0.10.2 | Up 3 months & 5 days
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 10 hours
  #  - www2 - 201.0.10.3 | Up 23 days
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 10 hours
  #  - www3 - 201.0.10.4 | Up 2 days
  #    * z918yd1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | X, Down 14 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 7 hours

// Cluster with two envirments
* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.10 | Up 3 months & 7 days
  #    * FedOps[x01] : "https://githib.com/wmiller/FedOps" | √, Built...1 minute ago
  #    * AmazingOtherApp[x02] : "https://githib.com/AmazingUser/AmazingOtherApp" | √, Built...2 hours ago
  #      dev | √, Pass...10 minutes | PUSH
  #        qzj1h8o -> trucks.www1-dev
  #      prod | √, Synced
  #        qzj1h8o -> trucks.www1, trucks.www2, trucks.www3
  #
  #  - cd-uswest2 - 201.0.10.11 | Up 3 months & 7 days
  #    * FedOps-www[x01] : "https://githib.com/wmiller/FedOps-www" | √, Built...2 hours ago
  #      dev | √, Pass...10 minutes | HOLD
  #        919jcah -> trucks.www1-dev
  #      prod | X, Not Synced
  #        s819af1 -> trucks.www1, trucks.www2, trucks.www3
  # trucks
  #  - www1-dev - 201.0.10.2 | Up 4 months & 9 days
  #    * 919jcah | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 20 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 10 hours
  #  - www1 - 201.0.10.2 | Up 3 months & 5 days
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 9 hours
  #  - www2 - 201.0.10.3 | Up 23 days
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 9 hours
  #  - www3 - 201.0.10.4 | Up 2 days
  #    * z918yd1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | X, Down 14 hours
  #    * qzj1h8o | AmazingOtherApp[x02] : warehouses.uswest1.AmazingOtherApp[x02] | √, Up 7 hours

* fedops ssh trucks.www1
  >$

* cake


-------------------------------------------

### Fresh Cluster with Manifest ###

* fedops init
  # enter provider
  > digital ocean
  # enter encryption key
  > fedops_key
  # enter digital ocean api key
  > api_key

* fedops container create some.service
  # enter encryption key
  >
* fedops use some.manifest
  # enter encryption key
  >

# wait a few seconds

* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.1 | Up 1 min
  #    * FedOps-www[x02] : "https://githib.com/wmiller/FedOps-www" | √, Built...2 hours ago
  #      s819af1 -> trucks.www1
  #
  # trucks
  #  - www1 - 201.0.10.2 | Up 1 min
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours

-------------------------------------------

### Fresh Cluster without Manifest ###

* fedops init
  # enter provider
  > digital ocean
  # enter encryption key
  > fedops_key
  # enter digital ocean api key
  > api_key

* fedops warehouse create [warehouse_name]
  # enter encryption key
  >
  # enter name
  > warehouse_name
* fedops truck create [truck_name]
  # enter encryption key
  >
  # enter name
  > truck_name

* fedops container create -w cd-uswest1 some.service
  # enter encryption key
  >

* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.1 | Up 1 min
  #    * FedOps-www[x02] : "https://githib.com/wmiller/FedOps-www" | _, Building...
  #
  # trucks
  #  - www1 - 201.0.10.2 | Up 1 min
  #

* fedops fork [truck_name] [forked_truck_name]
  # enter encryption key
  >
  # enter truck name
  > truck_name
  # enter fork name
  > forked_truck_name

* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.1 | Up 2 min
  #    * FedOps-www[x02] : "https://githib.com/wmiller/FedOps-www" | √, Built...1 min ago
  #
  # trucks
  #  - www1 - 201.0.10.2 | Up 2 min
  #  - www1-fork - 201.0.10.2 | Up 1 min
  #

* fedops env create [env_name]
  # enter encryption key
  >
  # enter env name
  > env_name

* fedops info // Example Output
  # enter encryption key
  >
  # warehouses
  #  - cd-uswest1 - 201.0.10.1 | Up 3 min
  #    * FedOps-www[x02] : "https://githib.com/wmiller/FedOps-www" | √, Built...2 hours ago
  #      prod | √, Pass...10 minutes
  #        s819af1 -> trucks.www1
  #
  # trucks
  #  - www1 - 201.0.10.2 | Up 3 min
  #    * s819af1 | FedOps-www[x01] : warehouses.uswest2.FedOps-www[x02] | √, Up 10 hours
  #  - www1-fork - 201.0.10.2 | Up 3 month & 5 days
  #

## About ##

## Liceneces ##
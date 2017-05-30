device-hub
-----------

[![Build Status](https://travis-ci.org/thingful/device-hub.svg?branch=master)](https://travis-ci.org/thingful/device-hub)

Transforms output from one or many IOT devices via one or many protocols into a common semantically understood format.
The logic to transform the data is via a "device profile". Device profiles are written in java-script.

![device-hub]( docs/device-hub-overview.png)

License
-------

Copyright © 2017 thingful

Released under the terms of "DECODE Accepted Software License"

Introduction
------------

device-hub is operated by installing a set of configured listeners, endpoint and device profiles.

device-hub is configured via the device-hub-cli.

device-hub stores its running configuration in a local boltdb database.

device-hub can run in insecure mode or using TLS in a mutual authentication model. See [security.md]( docs/security.md ) for details

Supported message formats

Transport               | Notes
------------------------|----------------------------------------------------------------
`CSV`                   |
`JSON`                  |
`RAW BYTES`             |

Supported listener transports

Transport               | Notes
------------------------|----------------------------------------------------------------
`HTTP`                  |
`MQTT`                  |

Supported endpoints

Transport               | Notes
------------------------|----------------------------------------------------------------
`STDOUT`                |
`HTTP`                  |

Example configuration files are in ./test-configurations/

The entity connecting a listener to a device profile and then an endpoint is called a 'pipe'.

On startup device-hub will restart all existing pipes.

Platforms
---------

device-hub supports all the platforms golang supports.

The Makefile contains build targets for linux amd64, raspberrypi arm and darwin (mac).

Build
-----

Install golang, docker (if you want to run the integration tests or test with a local mqtt server)

Get the code -

```
go get github.com/thingful/device-hub

```

The Makefile is documented

```
make help
```

Run the tests

```
make test
```

Build executables for all platforms

```
make all
```

Output is built to ./tmp/build/

Run
---

Start the device-hub

```
./device-hub-linux-amd64
```

Configure with the cli

To import a folder of configuration files

```
./device-hub-cli-linux-amd64 create -d=./test-configurations/
```

Files can also be imported on an individual basis

```
./device-hub-cli-linux-amd64 create -f=./test-configurations/mqtt_listener.yaml
```

The running configuration can be inspected

```
./device-hub-cli-linux-amd64 get all
```

Create some 'pipes' that listen via http on uri /a and /b and output to std output

```
./device-hub-cli-linux-amd64 start -e=stdout-endpoint -l=http-listener-local-port-8085 -u=/a thingful/device-1
./device-hub-cli-linux-amd64 start -e=stdout-endpoint -l=http-listener-local-port-8085 -u=/b thingful/device-2
```

It is possible to 'tag' messages with some user defined information when starting the 'pipe'

```
./device-hub-cli-linux-amd64 start -e=stdout-endpoint -l=http-listener-local-port-8085 -u=/c -t=foo:bar thingful/device-3

```

In the above example any messages received on '/c' will be tagged with the key value pair "foo" and "bar"


Send some messages 

```
curl -X POST -d '{"my-value": true}' 0.0.0.0:8085/a
curl -X POST -d '{"value": "22", "deviceId": "23", "createdAt": "1488205809000"}' 0.0.0.0:8085/b
```

Pipe through from MQTT -

Start the MQTT server

```
docker-compose up
```

```
./device-hub-cli-linux-amd64 start -e=stdout-endpoint -l=mqtt-listener-local-port-1883 -u=/some-mqtt-uri thingful/device-2
```

Send a message via mqtt to 0.0.0.0:1883 e.g. using MQTTLens (https://chrome.google.com/webstore/detail/mqttlens/hemojaaeigabkbcookmlgmdigohjobjm?hl=en)

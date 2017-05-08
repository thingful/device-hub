device-hub
-----------

Takes output from one or many IOT devices using various protocols (HTTP, MQTT) and expands it either transforming it into JSON format and/or adding in some metadata.
The logic to transform the data is via a "device profile". Device profiles are written in java-script.

```javascript
function decode (input) {
    console.log("decode called")
    return input
}

```

License
-------

Copyright © 2017 thingful
Released under the terms of "DECODE Accepted Software License"

Build
-----

Install golang

Get the code -

```
go get github.com/thingful/device-hub

```

Run the tests

```
make test
```

Build linux, mac and raspberry-pi executables

```
make all
```

Output is built to ./tmp/build/

Configuration
---------------

Example configuration files are in ./test-configurations/

Run
---

Start the device-hub 

```
./device-hub-linux-amd64
```

Configure with the cli

```
./device-hub-cli-linux-amd64 create -f=./test-configurations/http_listener.yaml
./device-hub-cli-linux-amd64 create -f=./test-configurations/mqtt_listener.yaml
./device-hub-cli-linux-amd64 create -f=./test-configurations/std_out_endpoint.yaml
./device-hub-cli-linux-amd64 create -f=./test-configurations/profile_script.yaml
./device-hub-cli-linux-amd64 create -f=./test-configurations/profile_script_transform.yaml
```

Create some 'pipes' that listen via http on uri /a and /b and out put to std output

```
./device-hub-cli-linux-amd64 start -e=mVZY2J0 -l=L1EYWOW -u=/a thingful/device-1
./device-hub-cli-linux-amd64 start -e=mVZY2J0 -l=L1EYWOW -u=/b thingful/device-2
```

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
./device-hub-cli-linux-amd64 start -e=mVZY2J0 -l=pjV4jN2 -u=/some-mqtt-uri thingful/device-2
```

Send a message via mqtt to 0.0.0.0:1883 e.g. using MQTTLens (https://chrome.google.com/webstore/detail/mqttlens/hemojaaeigabkbcookmlgmdigohjobjm?hl=en)

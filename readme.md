device-hub
-----------


Takes output from one or many IOT devices using various protocols (HTTP, MQTT) and expands it either transforming it into JSON format and/or adding in some metadata.
The logic to transform the data is written in java-script.

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

Javascript
----------

An example javascript file is in ./cmd/expando/decode.js 


Run
---

Pipe through from standard input -

```
 echo '{"value": "22", "deviceId": "23", "createdAt": "1488205809000"}' | ./device-hub -script="$(cat ./decode.js)" -in=std
{"@context":{"decode":"http://decode.xxx","m3-lite":"http://purl.org/iot/vocab/m3-lite#"},"@id":"decode:/23:1488205809000","@type":"m3-lite:AirPollutantSensor","createdAt":"1488205809000","deviceId":"23","domain":{"@type":"m3-lite:Environment"},"value":"22"}

```

Pipe through from MQTT -

Start the MQTT server

```
docker-compose up
```

```
./device-hub -in=mqtt -script="$(cat ./decode.js)"
```

Send a message via mqtt to 0.0.0.0:1883 e.g. using MQTTLens (https://chrome.google.com/webstore/detail/mqttlens/hemojaaeigabkbcookmlgmdigohjobjm?hl=en)

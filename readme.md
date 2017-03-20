expando
-------

Takes output from an IOT device and expands it either transforming it into JSON format and/or adding in some metadata. 
The logic to transform the data is written in javascript.

```javascript
function decode (input) {

    console.log("decode called")
    return input
}

```

Build
-----

Install golang

Get the code -

```
go get "bitbucket.org/tsetsova/decode-prototype

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
 echo '{"value": "22", "deviceId": "23", "createdAt": "1488205809000"}' | ./cmd/expando/expando -script="$(cat ./cmd/expando/decode.js)" -in=std
{"@context":{"decode":"http://decode.xxx","m3-lite":"http://purl.org/iot/vocab/m3-lite#"},"@id":"decode:/23:1488205809000","@type":"m3-lite:AirPollutantSensor","createdAt":"1488205809000","deviceId":"23","domain":{"@type":"m3-lite:Environment"},"value":"22"}

```

Pipe through from MQTT -

Start the MQTT server

```
docker-compose up
```

```
./cmd/expando/expando -in=mqtt -script="$(cat ./cmd/expando/decode.js)"
```

Send a message via mqtt to 0.0.0.0:1883 e.g. using MQTTLens (https://chrome.google.com/webstore/detail/mqttlens/hemojaaeigabkbcookmlgmdigohjobjm?hl=en)

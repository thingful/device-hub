Copyright © 2017 thingful

Released under the terms of "DECODE Accepted Software License"
<hr/>

Client CLI
=========================
```device-hub-cli [command]```


Client CLI Commands
=========================
```
- Name:           create
  Description:    Creates Listeners, Endpoints and Profile resources
  Flags:
  - Flag:         -d
    Large Format: --request-dir <string>
    Is required:  Yes if -f isn't specified
    Description:  Directory containing client request file(s) (must be json, yaml, or xml) 
    Parameter:    A filesystem path
    Example:      create -d=./test-configurations/

  - Flag:         -f
    Large Format: --request-file <string>
    Is required:  Yes if -d isn't specified
    Description:  Client request file (must be json, yaml, or xml); use "-" for stdin + json
    Parameter:    Directory containing client request file(s) (must be json, yaml, or xml) 
    Example:      create -f=./test-configurations/mqtt_listener.yaml


- Name:           delete
  Description:    Delete listener, profile and endpoint resources
  Flags:
  - Flag:         -d
    Large Format: --request-dir <string>
    Is required:  Yes if -f isn't specified
    Description:  Directory containing client request file(s) (must be json, yaml, or xml) 
    Parameter:    A filesystem path
    Example:      create -d=./test-configurations/

  - Flag:         -f
    Large Format: --request-file <string>
    Is required:  Yes if -d isn't specified
    Description:  Client request file (must be json, yaml, or xml); use "-" for stdin + json
    Parameter:    Directory containing client request file(s) (must be json, yaml, or xml) 
    Example:      create -f=./test-configurations/mqtt_listener.yaml

- Name:           describe
  Description:    Describe parameters for endpoint and listeners

- Name:           get
  Description:    Display one or many resources

- Name:           list
  Description:    List running pipes

- Name:           start
  Description:    Start processing messages on an URI

- Name:           stop
  Description:    Stop processing messages on an URI

- Name:           version
  Description:    Display version information

- Name:           help
  Description:    Help about any command
```
CLI Global Flags
=================
```
- Flag:           -s
  Large Format:   --server-addr <string>
  Description:    Server address in form of host:port
  Default:        127.0.0.1:50051

- Flag:           -o
  Large Format:   --response-format <string>
  Description:    response format (json, prettyjson, yaml, or xml)
  Default:        json

- Flag:           -p
  Large Format:   --print-sample-request
  Description:    Print sample request file and exit

- Flag:           --timeout <duration>
  Description:     Client connection timeout
  Default:          10s
```
CLI Global Security Flags
=========================
```
- Flag:           --auth-token <string>
  Description:    Authorization token

- Flag:           --auth-token-type <string>
  Description:    Authorization token type
  Default:        Bearer

- Flag:           --jwt-key <string>
  Description:    Jwt key

- Flag:           --jwt-key-file <string>
  Description:    Jwt key file
```
CLI Global TLS Security Flags
=============================
```
- Flag:           --tls
  Description:    Enable TLS

- Flag:           --tls-ca-cert-file <string>
  Description:    CA certificate file

- Flag:           --tls-cert-file <string>
  Description:    Client certificate file

- Flag:           --tls-insecure-skip-verify
  Description:    INSECURE: Skip TLS checks

- Flag:           --tls-key-file <string>
  Description:    Client key file

- Flag:           --tls-server-name <string>
  Description:    TLS Server name override
  ```
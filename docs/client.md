Copyright © 2017 thingful

Released under the terms of "DECODE Accepted Software License"
<hr/>

Client CLI
=========================
```
device-hub-cli [command]
device-hub-cli [global-flags] [command]
```


Client CLI Commands
=========================
```
- Command Name:   create [-d|-f]=<string>
  Description:    Creates Listeners, Endpoints, Profile resources and start process if -d flag is used
  Flags:
  - Flag:         -d
    Large Format: --request-dir <string>
    Is required:  Yes if -f isn't specified
    Description:  Directory containing client yaml request file(s) 
    Parameter:    A filesystem path with valid yaml resources and processes files
    Example:      device-hub-cli create -d=./test-configurations/

  - Flag:         -f
    Large Format: --request-file <string>
    Is required:  Yes if -d isn't specified
    Description:  Client yaml request file
    Parameter:    Path containing client request file(s) (must be yaml) 
    Example:      device-hub-cli create -f=./test-configurations/mqtt_listener.yaml


- Command Name:   delete [-d|-f]=<string>
  Description:    Delete listener, profile, endpoint resources and stop processes if process files are specified
  Flags:
  - Flag:         -d
    Large Format: --request-dir <string>
    Is required:  Yes if -f isn't specified
    Description:  Directory containing client yaml request file(s)
    Parameter:    A filesystem path
    Example:      device-hub-cli delete -d=./test-configurations/

  - Flag:         -f
    Large Format: --request-file <string>
    Is required:  Yes if -d isn't specified
    Description:  Client request yaml file
    Parameter:    Directory containing client request file(s) 
    Example:      device-hub-cli delete -f=./test-configurations/mqtt_listener.yaml


- Command Name:   describe [listener|endpoint] [mqtt|stdout]
  Description:    Describe parameters for endpoints and listeners
  Example:        device-hub-cli listener mqtt


- Command Name:   show [listener|endpoint|profile|all]
  Description:    Display one or many resources by type or using "all" as * filter, 
  Example:        device-hub-cli show all 


- Command Name:   status
  Description:    List running pipes


- Command Name:   start [-f <file>] [-e <string> -l <string> -u <string>] <string>
  Description:    Start processing messages on an URI
  Flags:
  - Flag:         -f <file>
    Description:  Start processing messages using a yaml config file
    Large Format: --request-file <file>
    Is required:  No

  - Flag:         -e
    Description:  Endpoint uid to push messages to, may be specified multiple times
    Large Format: --endpoint <stringSlice>
    Is required:  Yes if -f isn't specified

  - Flag:         -l
    Description:  Listener uid to accept messages on
    Large Format: --listener <string>
    Is required:  Yes if -f isn't specified
  
  - Flag:         -u
    Description:  Uri to listen on
    Large Format: --uri <string>
    Is required:  Yes if -f isn't specified

  - Flag:         -t
    Description:  Colon separated (k:v) runtime tags to attach to requests, may be specified multiple times
    Large Format: --tags <stringSlice>
    Is required:  Yes if -f isn't specified
  
  Example:        device-hub-cli start -e=stdout-endpoint -l=http-listener-local-port-8085 -u=/a thingful/helsinki-bus


- Command Name:   stop <string>
  Description:    Stop processing messages on an URI
  Example:        device-hub-cli stop /a


- Command Name:   version
  Description:    Display version information
  Example:        device-hub-cli version


- Command Name:   help <string>
  Description:    Help about any command
  Example:        device-hub-cli help stop
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

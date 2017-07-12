Copyright © 2017 thingful

Released under the terms of "DECODE Accepted Software License"
<hr/>

Server CLI
=========================
```device-hub [command]```


CLI Commands
=========================
```
- Name:           server
  Description:    Start device hub

- Name:           version
  Description:    Display version information

- Name:           help
  Description:    Help about any command
```
CLI Global Flags
=================
```
- Flag:           -c
  Large Format:   --config-file
  Description:    Enable config file, fields in config file will override flags and env vars
  Default:        false

- Flag:           --config-path
  Description:    Path to config file
  Default:        ./config.yaml

- Flag:           -b
  Large Format:   --binding <string>
  Description:    Binding address in form of {ip}:port
  Default:        :50051

- Flag:           --data-dir <string>
  Description:    Path to db folder
  Default:        Current directory "."

- Flag:           --data-impl <string>
  Description:    datastore implementation to use, valid values are 'boltdb' or 'filestore'
  Default:        boltdb

- Flag:           --log-file
  Description:    enable log to file
  Default:        false

- Flag:           --log-path <string>
  Description:    path to log file
  Default:        ./device-hub.log

- Flag:           --log-syslog
  Description:    enable log to local SYSLOG
  Default:        false

```
CLI Global TLS Security Flags
=============================
```
- Flag:           --tls
  Description:    Enable tls

- Flag:           --tls-ca-cert-file <string>
  Description:    CA certificate file

- Flag:           --tls-cert-file <string>
  Description:    Client certificate file

- Flag:           --tls-key-file <string>
  Description:    Client key file
```

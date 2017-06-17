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
- Flag:           -b
  Large Format:   --binding <string>
  Description:    Binding address in form of {ip}:port
  Default:        :50051

- Flag:           --data <string>
  Description:    Path to db folder
  Default:        Current directory "."
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
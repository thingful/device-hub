Security & Authentication
=========================

device-hub uses `Transport Layer Security <https://en.wikipedia.org/wiki/Transport_Layer_Security>`_ (TLS) to authenticate components, in a `mutual authentication <https://en.wikipedia.org/wiki/Mutual_authentication>`_ model.

This ensures that all communicating parties are communicating with a verified component, helping to prevent unauthorized access, and mitigating some potential attack vectors.

Mutual Authentication Overview
==============================

The certificates and private keys for the device-hub service are installed alongside the master root public certificate file.

device-hub's API end users are issued their own certificate and private key, also with a copy of the master root public certificate file.

This allows all components to establish both a private channel of communication and a means of verifying identity; the client validates the server certificate was signed by the certifying authority, while the server mutually verifies the client's certificate was signed by the same authority.

Security Benefits
=================

The TLS client certification layer used by device-hub provides a number of security benefits.

- Prevents unauthorized requests to the device-hub control API.
- Prevents unauthorized connections to the device-hub control API.
- Encrypts communications between all components.

Risks
=====

device-hub's authentication layer does not completely guarantee the security of an installation; it relies on the private keys of the installations certificates being kept secret.

For example, if a malicious user were able to gain root SSH access to the machine running the device-hub client, they would be able to copy the client's private key and therefore be able to administer a remote device-hub deployment.

It is therefore very important that you ensure the private keys are kept secure; they should not be copied or shared insecurely.
When copying certificates and private keys as part of the device-hub installation process, the files must be copied using a secure and encrypted transfer medium such as SSH, SCP or SFTP.

Other measures that would normally be taken to secure a server should still be implemented; device-hub is not a full server stack and its security layer does not prevent your server being hacked, rather it mitigates the likelihood of device-hub services being used as a vector to do so.

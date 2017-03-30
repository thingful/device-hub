package endpoint

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	hub "github.com/thingful/device-hub"
	multichainClient "github.com/thingful/tm/bc/multichain"
	"github.com/thingful/tm/crypto"
	"github.com/thingful/tm/encoders"
)

func NewMultichainEndpoint(publicKey *rsa.PublicKey, deviceID string, notifyURI string) multichain {
	return multichain{
		publicKey: publicKey,
		deviceID:  deviceID,
		notifyURI: notifyURI,
	}
}

type multichain struct {
	publicKey *rsa.PublicKey
	deviceID  string
	notifyURI string
}

func (m multichain) Write(message hub.Message) error {

	// TODO : pass in client details
	client := multichainClient.New("chain1",
		"192.168.56.33",
		4800,
		"multichainrpc",
		"8VY5wGdNGuKQKYnmtEss9gPxpFLzAk6mkgSQEKvZ7Pce")

	// create one time password (OTP) to encrypt the data with
	otp := make([]byte, 32)
	_, err := rand.Read(otp)

	if err != nil {
		return err
	}

	dataBytes, err := encoders.JSONMarshaler{}.Marshal(message)

	if err != nil {
		return err
	}

	encryptedData, err := crypto.AESEncypt(otp, dataBytes)

	if err != nil {
		return err
	}

	// write data to the blockchain using the OTP
	dataStream := m.deviceID

	client.Create(dataStream, false)
	transactionID, err := client.PublishToStream(dataStream, "data", encryptedData)

	if err != nil {
		return err
	}

	// send notification to the entitlement owners data that there is new data
	client.Create(m.notifyURI, true)

	accessTokenMessage := multichainClient.AccessTokenMessage{
		Secret: otp,
		TxnID:  transactionID,
	}

	notifyEncoder := encoders.New(
		encoders.JSONMarshaler{},
		encoders.RSAPublicKeyEncrypter(m.publicKey))

	j, err := notifyEncoder.Encode(accessTokenMessage)

	if err != nil {
		return err
	}

	// send users update of new information
	response, err := client.PublishToStream(m.notifyURI, "access-token", j)

	if err != nil {
		return err
	}

	fmt.Println("published access details : ", response)

	return nil
}

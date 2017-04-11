// Copyright © 2017 thingful

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var RootCmd = &cobra.Command{
	Use: "device-hub-cli",
}

func init() {

	var hubAddress string

	// Client can run either in insecure mode or provide details for mutual tls
	// The default is for secure connections to be used.
	var insecure bool

	// if insecure == false the following need to be set
	var certFilePath string
	var keyFilePath string
	var trustedCAFilePath string

	RootCmd.PersistentFlags().StringVarP(&hubAddress, "binding", "b", "localhost:50051", "RPC binding for the device-hub daemon.")

	RootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "Switch off Mutual TLS authentication.")

	RootCmd.PersistentFlags().StringVar(&certFilePath, "cert-file", "", "Certificate used for SSL/TLS RPC connections to the device-hub daemon.")
	RootCmd.PersistentFlags().StringVar(&keyFilePath, "key-file", "", "Key file for the certificate (--cert-file).")
	RootCmd.PersistentFlags().StringVar(&trustedCAFilePath, "trusted-ca-file", "", "Trusted certificate authority.")

	versionCommand := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(hub.ClientVersionString())
			return nil
		},
	}

	RootCmd.AddCommand(versionCommand)

	pipeCommands := &cobra.Command{
		Use: "pipe",
	}

	listCommand := &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := NewGRPCConnection(insecure, hubAddress, certFilePath, keyFilePath, trustedCAFilePath)

			if err != nil {
				return fmt.Errorf("did not connect: %v", err)
			}

			defer conn.Close()
			c := proto.NewHubClient(conn)

			r, err := c.PipeList(context.Background(), &proto.PipeListRequest{})

			if err != nil {
				return err
			}

			for _, p := range r.Pipes {
				// TODO : find a really good terminal library!
				fmt.Printf("URI: '%s' \tSTATUS: %s\n", p.Uri, p.State.String())

				endpointsSummary := []string{}
				for _, e := range p.Endpoints {
					endpointsSummary = append(endpointsSummary, fmt.Sprintf("%s [%s]", e.Uid, e.Type))
				}

				fmt.Printf("%s [%s] -> %s [%s] -> %s\n", p.Listener.Uid, p.Listener.Type, p.Profile.Name, p.Profile.Version, strings.Join(endpointsSummary, " && "))
				fmt.Printf("Ok: %d Errors: %d Total : %d\n", p.MessageStats.Ok, p.MessageStats.Errors, p.MessageStats.Total)

			}

			return nil
		}}

	var uri string

	deleteCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {

			conn, err := NewGRPCConnection(insecure, hubAddress, certFilePath, keyFilePath, trustedCAFilePath)

			if err != nil {
				return fmt.Errorf("did not connect: %v", err)
			}

			defer conn.Close()
			c := proto.NewHubClient(conn)

			r, err := c.PipeDelete(context.Background(), &proto.PipeDeleteRequest{
				Uri: uri,
			})

			if err != nil {
				return err
			}

			fmt.Println(r)

			return nil
		}}

	deleteCmd.Flags().StringVar(&uri, "uri", "", "Uri of pipe to delete")

	pipeCommands.AddCommand(listCommand, deleteCmd)
	RootCmd.AddCommand(pipeCommands)

}

func main() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func NewGRPCConnection(insecure bool, address string, certFilePath, keyFilePath, trustedCAFilePath string) (*grpc.ClientConn, error) {

	if insecure {
		return insecureGRPCConnection(address)
	}
	return secureGRPCConnection(address, certFilePath, keyFilePath, trustedCAFilePath)
}

func insecureGRPCConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithInsecure())
}

func secureGRPCConnection(address string, certFilePath, keyFilePath, trustedCAFilePath string) (*grpc.ClientConn, error) {

	certificate, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not load client key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(trustedCAFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not read ca certificate: %s", err)
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, errors.New("failed to append ca certs")
	}

	// Parse out the serverName
	urz, err := url.Parse(address)

	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   urz.Host,
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	// Create a connection with the TLS credentials
	return grpc.Dial(address, grpc.WithTransportCredentials(creds))
}

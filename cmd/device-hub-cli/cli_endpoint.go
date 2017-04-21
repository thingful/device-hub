// Copyright © 2017 thingful

package main

/*
var endpointCommand = &cobra.Command{
	Use:   "endpoint",
	Short: "Add, Delete and List endpoints.",
}

var endpointAddCommand = &cobra.Command{
	Use: "add",
	Long: `Add a new endpoint

You can use environment variables with the same name of the command flags.
All caps and s/-/_, e.g. SERVER_ADDR.`,
	Example: `
Save a sample request to a file (or refer to your protobuf descriptor to create one):
	device-hub-cli endpoint add -p > req.json
Submit request using file:
	device-hub-cli endpoint add -f req.json`,
	Run: func(cmd *cobra.Command, args []string) {
		v := proto.EndpointAddRequest{
			Endpoint: &proto.Endpoint{},
		}
		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.EndpointAdd(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}*/

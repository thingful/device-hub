// Copyright © 2017 thingful

package main

/*
var pipeCommand = &cobra.Command{
	Use:   "pipe",
	Short: "Add, Delete and List pipes.",
}

var pipeListCommand = &cobra.Command{
	Use: "list",
	Long: `List pipes

You can use environment variables with the same name of the command flags.
All caps and s/-/_, e.g. SERVER_ADDR.`,
	Example: `
Save a sample request to a file (or refer to your protobuf descriptor to create one):
	device-hub-cli pipe list -p > req.json
Submit request using file:
	device-hub-cli list -f req.json`,
	Run: func(cmd *cobra.Command, args []string) {
		var v proto.PipeListRequest
		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.PipeList(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

var pipeAddCommand = &cobra.Command{
	Use: "add",
	Long: `Add a new pipe

You can use environment variables with the same name of the command flags.
All caps and s/-/_, e.g. SERVER_ADDR.`,
	Example: `
Save a sample request to a file (or refer to your protobuf descriptor to create one):
	device-hub-cli pipe add -p > req.json
Submit request using file:
	device-hub-cli pipe add -f req.json`,
	Run: func(cmd *cobra.Command, args []string) {
		var v proto.PipeAddRequest
		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.PipeAdd(context.Background(), &v)

			if err != nil {
				return err
			}

			return out.Encode(resp)

		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

var pipeDeleteCommand = &cobra.Command{
	Use: "delete",
	Long: `Delete an existing pipe

You can use environment variables with the same name of the command flags.
All caps and s/-/_, e.g. SERVER_ADDR.`,
	Example: `
Save a sample request to a file (or refer to your protobuf descriptor to create one):
	device-hub-cli pipe delete -p > req.json
Submit request using file:
	device-hub-cli delete -f req.json`,
	Run: func(cmd *cobra.Command, args []string) {
		var v proto.PipeDeleteRequest
		err := roundTrip(v, func(cli proto.HubClient, in iocodec.Decoder, out iocodec.Encoder) error {

			err := in.Decode(&v)
			if err != nil {
				return err
			}

			resp, err := cli.PipeDelete(context.Background(), &v)

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

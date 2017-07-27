// Copyright © 2017 thingful

package runtime

import (
	"context"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
	"github.com/thingful/device-hub/utils"
)

// loop orchestrates the managed runtime loop
func loop(ctx context.Context,
	p *Pipe,
	listener hub.Listener,
	endpoints map[string]hub.Endpoint,
	channel hub.Channel,
	logger utils.Logger,
	tags map[string]string,
	options map[string]interface{}) {

	scripter := engine.New(logger)
	if options["geolocation"].(bool) {
		scripter.SetGeoLocation(options["lat"].(float64), options["lng"].(float64))
	}

	// ensure the map for the Statistics.Sent is set up correctly
	for k, _ := range endpoints {
		p.Statistics.Sent[k] = &proto.Counters{}
	}
	for {

		select {
		case <-ctx.Done():
			p.State = proto.Pipe_STOPPED
			err := channel.Close()

			if err != nil {
				logger.Error(err)
			}

			return

		case err := <-channel.Errors():

			p.Statistics.Received.Total++
			p.Statistics.Received.Errors++
			logger.Error(err)

		case input := <-channel.Out():
			p.Statistics.Received.Total++
			p.Statistics.Received.Ok++

			input.Metadata[hub.ENGINE_TIMESTAMP_START_KEY] = time.Now().UTC()

			output, err := scripter.Execute(p.Profile.Script, input)

			output.Metadata[hub.ENGINE_TIMESTAMP_END_KEY] = time.Now().UTC()

			p.Statistics.Processed.Total++

			if err != nil {
				p.Statistics.Processed.Errors++

				output.Metadata[hub.ENGINE_OK_KEY] = false
				output.Metadata[hub.ENGINE_ERROR_KEY] = err.Error()

				logger.Error(err)
			} else {

				output.Metadata[hub.ENGINE_OK_KEY] = true

				p.Statistics.Processed.Ok++
			}

			output.Metadata[hub.PROFILE_NAME_KEY] = p.Profile.Name
			output.Metadata[hub.PROFILE_VERSION_KEY] = p.Profile.Version
			output.Metadata[hub.RUNTIME_VERSION_KEY] = hub.SourceVersion
			output.Tags = tags

			output.Schema = p.Profile.Schema

			for k, _ := range endpoints {

				p.Statistics.Sent[k].Total++

				err = endpoints[k].Write(output)

				if err != nil {
					p.Statistics.Sent[k].Errors++
					logger.Error(err)
				} else {
					p.Statistics.Sent[k].Ok++
				}
			}
		}
	}
}

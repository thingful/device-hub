// Copyright © 2017 thingful

package runtime

import (
	"context"
	"log"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/engine"
	"github.com/thingful/device-hub/proto"
)

// loop orchestrates the managed runtime loop
func loop(ctx context.Context,
	p *Pipe,
	listener hub.Listener,
	endpoints map[string]hub.Endpoint,
	channel hub.Channel,
	tags map[string]string) {

	scripter := engine.New()

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
				log.Println(err)
			}

			return

		case err := <-channel.Errors():

			p.Statistics.Received.Total++
			p.Statistics.Received.Errors++
			log.Println(err)

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

				log.Println(err)
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

				// TODO : do something more useful with this error
				err = endpoints[k].Write(output)

				if err != nil {
					p.Statistics.Sent[k].Errors++
					log.Println(err)
				} else {
					p.Statistics.Sent[k].Ok++
				}
			}
		}
	}
}

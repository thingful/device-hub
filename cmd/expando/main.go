package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/thingful/expando"
	"github.com/thingful/expando/engine"
)

func main() {

	var in string
	var scriptContents string

	flag.StringVar(&in, "in", "{}", "json input")
	flag.StringVar(&scriptContents, "script", "function decode( input ){ return input }", "js to transform input")

	/*	j := `{
			"value": "30",
			"deviceId": "23",
			"createdAt": "1487941771000"
		}`
	*/
	flag.Parse()

	scripter := engine.New()
	input := expando.Input{Payload: []byte(in)}

	script := expando.Script{
		Runtime:  expando.Javascript,
		Input:    expando.JSON,
		Contents: scriptContents,

		/*`function decode (input) {

					// define the terms against the m3-lite ontology
					// http://ontology.fiesta-iot.eu/ontologyDocs/fiesta-iot/doc
					input['@context'] = {
		                'm3-lite': 'http://purl.org/iot/vocab/m3-lite#'
					}

					// it is an air pollutant sensor
					input['@type'] = "m3-lite:AirPollutantSensor"

					// environment based
					input['domain'] = {
		                "@type" : "m3-lite:Environment"
		            }

					// TODO : what is the value, unit?

					return input
					}`,*/
	}

	output, err := scripter.Execute(input, script)
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(output, "", "   ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

}

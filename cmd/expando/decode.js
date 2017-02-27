
// decode adds ontology information to the input
function decode (input) {

    // you can log to debug your script!
    //console.log("decode called")

    input['@context'] = {
        // define the terms against the m3-lite ontology
        // http://ontology.fiesta-iot.eu/ontologyDocs/fiesta-iot/doc
        'm3-lite': 'http://purl.org/iot/vocab/m3-lite#',
        // add in a decode context for the unique id
        'decode' : 'http://decode.xxx',
    }

    input['@id'] = "decode:/" + input['deviceId'] + ':' + input['createdAt']

    // it is an air pollutant sensor
    input['@type'] = "m3-lite:AirPollutantSensor"

    // in the environment domain
    input['domain'] = {
    "@type" : "m3-lite:Environment"
    }

    // TODO : what is the value, data format, unit?
    return input
}

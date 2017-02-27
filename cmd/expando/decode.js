
// decode will add ontology information to the input
function decode (input) {
    // you can log to debug your script!
    //console.log("decode called")

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

    // TODO : what is the value, data format, unit?
    return input
}

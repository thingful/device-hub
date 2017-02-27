package expando

type InputType string

type Runtime int

const (
	Javascript Runtime = 1

	Raw  InputType = "raw"
	CSV  InputType = "csv"
	XML  InputType = "xml"
	JSON InputType = "json"
)

type Script struct {
	Main     string
	Runtime  Runtime
	Input    InputType
	Contents string
	Metadata map[string]interface{}
}

type Input struct {
	Payload  []byte
	Metadata map[string]interface{}
}

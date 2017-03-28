package expando

type InputType string

type Runtime string

const (
	Javascript Runtime = "javascript"

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

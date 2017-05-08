// Copyright © 2017 thingful

package engine

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
	Main     string                 `json:"main"`
	Runtime  Runtime                `json:"runtime"`
	Input    InputType              `json:"input"`
	Contents string                 `json:"contents"`
	Metadata map[string]interface{} `json:"metadata"`
}

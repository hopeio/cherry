package encoding

type Format string

const (
	JSON     Format = "json"
	YAML     Format = "yaml"
	TOML     Format = "toml"
	YML      Format = "yml"
	Protobuf Format = "protobuf"
	XML      Format = "xml"
)

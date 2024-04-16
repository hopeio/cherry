package pretty

import "github.com/tidwall/pretty"

func Pretty(json []byte) []byte { return pretty.PrettyOptions(json, nil) }

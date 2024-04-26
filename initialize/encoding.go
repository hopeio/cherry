package initialize

type Encoder interface {
	Encode(format string, v map[string]any) ([]byte, error)
}

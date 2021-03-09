package tool

import (
	"encoding/json"
	"io"
)

type JsonParse struct {
}

func Decoder(io io.ReadCloser, v interface{}) error {
	return json.NewDecoder(io).Decode(v)
}

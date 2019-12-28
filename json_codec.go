package cache

import (
	"encoding/json"
)

type jsonCodec struct{}

func (jc jsonCodec) Encode(val interface{}) ([]byte, error) {
	return json.Marshal(val)
}

func (jc jsonCodec) Decode(bytes []byte, dest interface{}) error {
	return json.Unmarshal(bytes, dest)
}

func JsonCodec() CodecInterface {
	return jsonCodec{}
}

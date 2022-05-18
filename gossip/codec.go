package gossip

import b64 "encoding/base64"

type MessageCodec interface {
	Encode(payload []byte) []byte
	Decode(payload []byte) ([]byte, error)
}

type B64Codec struct{}

func (c *B64Codec) Encode(payload []byte) []byte {
	enc := b64.StdEncoding
	encoded := make([]byte, enc.EncodedLen(len(payload)))
	enc.Encode(encoded, payload)
	return encoded
}

func (c *B64Codec) Decode(payload []byte) ([]byte, error) {
	enc := b64.StdEncoding
	decoded := make([]byte, enc.DecodedLen(len(payload)))
	n, err := enc.Decode(decoded, payload)
	if err != nil {
		return nil, err
	}
	return decoded[:n], nil
}

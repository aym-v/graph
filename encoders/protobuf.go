package encoders

import (
	"errors"

	"google.golang.org/protobuf/proto"
)

// Protobuf represents a Protocol Buffer encoder.
type Protobuf struct{}

// Map names
const (
	ProtobufEncoder = "protobuf"
)

// Errors
var (
	ErrInvalidProtoMsgEncode = errors.New("graph: Invalid protobuf proto.Message object passed to encode")
	ErrInvalidProtoMsgDecode = errors.New("graph: Invalid protobuf proto.Message object passed to decode")
)

// NewProtobuf returns a new Protobuf encoder.
func NewProtobuf() *Protobuf {
	return &Protobuf{}
}

// Encode encodes a protobuf message.
func (pb *Protobuf) Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	i, found := v.(proto.Message)
	if !found {
		return nil, ErrInvalidProtoMsgEncode
	}

	b, err := proto.Marshal(i)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Decode decodes a protobuf message.
func (pb *Protobuf) Decode(data []byte, vPtr interface{}) error {
	if _, ok := vPtr.(*interface{}); ok {
		return nil
	}
	i, found := vPtr.(proto.Message)
	if !found {
		return ErrInvalidProtoMsgDecode
	}

	return proto.Unmarshal(data, i)
}

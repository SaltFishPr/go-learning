package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const magicNumber byte = 42

var (
	// ErrMetaKVMissing some keys or values are missing.
	ErrMetaKVMissing = errors.New("wrong metadata lines. some keys or values are missing")
	// ErrMessageTooLong message is too long
	ErrMessageTooLong = errors.New("message is too long")
)

// CompressType defines decompression type.
type CompressType byte

const (
	// CompressTypeNone does not compress.
	CompressTypeNone CompressType = iota
	// CompressTypeGzip uses gzip compression.
	CompressTypeGzip
)

func (t CompressType) String() string {
	switch t {
	case CompressTypeNone:
		return "none"
	case CompressTypeGzip:
		return "gzip"
	default:
		return "unknown"
	}
}

// CodecType defines serialization type of payload.
type CodecType byte

const (
	// CodecTypeNone uses raw []byte and don't serialize/deserialize
	CodecTypeNone CodecType = iota
	// CodecTypeJSON for payload.
	CodecTypeJSON
	// CodecTypeProtoBuffer for payload.
	CodecTypeProtoBuffer
)

func (t CodecType) String() string {
	switch t {
	case CodecTypeNone:
		return "none"
	case CodecTypeJSON:
		return "json"
	case CodecTypeProtoBuffer:
		return "protobuf"
	default:
		return "unknown"
	}
}

type Header [12]byte

func (h Header) CheckMagicNumber() bool {
	return h[0] == magicNumber
}

func (h Header) Version() byte {
	return h[1]
}

func (h Header) CompressType() CompressType {
	return CompressType(h[2] >> 4)
}

func (h Header) CodecType() CodecType {
	return CodecType(h[2] & 0x0f)
}

// Message .
type Message struct {
	Header        *Header
	ServiceMethod string // format "Service.Method"
	Metadata      map[string]string
	Payload       []byte
}

func NewMessage() *Message {
	header := Header([12]byte{})
	header[0] = magicNumber
	return &Message{
		Header: &header,
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("Header: %+v, ServiceMethod: %s, Metadata: %+v, Payload: %s", m.Header, m.ServiceMethod, m.Metadata, string(m.Payload))
}

func (m *Message) Decode(r io.Reader) error {
	_, err := io.ReadFull(r, m.Header[:1])
	if err != nil {
		return fmt.Errorf("read magic number error: %w", err)
	}
	if !m.Header.CheckMagicNumber() {
		return fmt.Errorf("wrong magic number: %v", m.Header[0])
	}
	_, err = io.ReadFull(r, m.Header[1:])
	if err != nil {
		return fmt.Errorf("read header error: %w", err)
	}

	totalSizeData := make([]byte, 4)
	_, err = io.ReadFull(r, totalSizeData)
	if err != nil {
		return fmt.Errorf("read total size error: %w", err)
	}
	l := binary.BigEndian.Uint32(totalSizeData)

	data := make([]byte, int(l))
	_, err = io.ReadFull(r, data)
	if err != nil {
		return fmt.Errorf("read remaining error: %w", err)
	}

	n := 0
	// parse ServiceMethod
	l = binary.BigEndian.Uint32(data[n:4])
	n = n + 4
	nEnd := n + int(l)
	m.ServiceMethod = string(data[n:nEnd])
	n = nEnd

	// parse Metadata
	l = binary.BigEndian.Uint32(data[n : n+4])
	n = n + 4
	nEnd = n + int(l)
	if l > 0 {
		m.Metadata, err = decodeMetadata(l, data[n:nEnd])
		if err != nil {
			return err
		}
	}
	n = nEnd

	// parse Payload
	l = binary.BigEndian.Uint32(data[n : n+4])
	_ = l
	n = n + 4
	m.Payload = data[n:]

	return nil
}

func encodeMetadata(m map[string]string, bb *bytes.Buffer) {
	if len(m) == 0 {
		return
	}
	d := make([]byte, 4)
	for k, v := range m {
		binary.BigEndian.PutUint32(d, uint32(len(k)))
		bb.Write(d)
		bb.Write([]byte(k))
		binary.BigEndian.PutUint32(d, uint32(len(v)))
		bb.Write(d)
		bb.Write([]byte(v))
	}
}

func decodeMetadata(l uint32, data []byte) (map[string]string, error) {
	m := make(map[string]string, 10)
	var n uint32
	for n < l {
		// parse one key and value
		// key
		sl := binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l-4 {
			return m, ErrMetaKVMissing
		}
		k := string(data[n : n+sl])
		n = n + sl

		// value
		sl = binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l {
			return m, ErrMetaKVMissing
		}
		v := string(data[n : n+sl])
		n = n + sl
		m[k] = v
	}
	return m, nil
}

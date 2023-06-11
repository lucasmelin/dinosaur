package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	Id             uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func ParseHeader(header *bytes.Reader) Header {
	h := Header{}
	binary.Read(header, binary.BigEndian, &h.Id)
	binary.Read(header, binary.BigEndian, &h.Flags)
	binary.Read(header, binary.BigEndian, &h.NumQuestions)
	binary.Read(header, binary.BigEndian, &h.NumAnswers)
	binary.Read(header, binary.BigEndian, &h.NumAuthorities)
	binary.Read(header, binary.BigEndian, &h.NumAdditionals)

	return h
}

// ToBytes encodes a Header as bytes.
// In network packets, integers are always encoded in big endian.
func (h *Header) ToBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, h.Id)
	_ = binary.Write(buf, binary.BigEndian, h.Flags)
	_ = binary.Write(buf, binary.BigEndian, h.NumQuestions)
	_ = binary.Write(buf, binary.BigEndian, h.NumAnswers)
	_ = binary.Write(buf, binary.BigEndian, h.NumAuthorities)
	_ = binary.Write(buf, binary.BigEndian, h.NumAdditionals)
	return buf.Bytes()
}

func (h Header) String() string {
	return fmt.Sprintf(
		`Header{
    Id: %d,
    Flags: %d,
    NumQuestions: %d,
    NumAnswers: %d,
    NumAuthorities: %d,
    NumAdditionals: %d
  }`,
		h.Id,
		h.Flags,
		h.NumQuestions,
		h.NumAnswers,
		h.NumAuthorities,
		h.NumAdditionals,
	)
}

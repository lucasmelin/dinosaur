package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	// ID is a 16 bit identifier assigned by the program that generated the query.
	// This identifier is copied to the corresponding reply,
	// and can be used by the requester to match up replies to outstanding queries.
	ID uint16
	// Flags represents the following fields:
	//  - QR - a one bit field that specifies if this message is a query (0) or a response (1).
	//  - OPCODE - a four bit field that specifies the kind of query in this message. This value
	//    is set by the originator of a query and copied into the response. The values are a standard query (0),
	//    an inverse query (1), a server status request (2), reserved for future use (3-15).
	//  - AA - a one bit field that is valid in responses and specifies that the responding name server is an authority
	//    for the domain name in the question. Note that the contents of the answer section may have multiple owner names
	//    because of aliases. The AA bit corresponds to the name which matches the query name, or the fist owner name
	//    in the answer section.
	//  - TC - a one bit field that specifies that this message was truncated due to length greater than that permitted
	//    on the transmission channel.
	//  - RD - a one bit field that may be set in a query and is copied into the response. If recursion desired is set,
	//    it directs the name server to pursue the query recursively. Recursive query support is optional.
	//  - RA - a one bit field that is set or cleared in a response, and denotes whether recursive query support is
	//    available in the name server.
	//  - Z - a one bit field that is reserved for future use. Must be zero in all queries and responses.
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func ParseHeader(header *bytes.Reader) Header {
	h := Header{}
	binary.Read(header, binary.BigEndian, &h.ID)
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
	_ = binary.Write(buf, binary.BigEndian, h.ID)
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
    ID: %d,
    Flags: %d,
    NumQuestions: %d,
    NumAnswers: %d,
    NumAuthorities: %d,
    NumAdditionals: %d
  }`,
		h.ID,
		h.Flags,
		h.NumQuestions,
		h.NumAnswers,
		h.NumAuthorities,
		h.NumAdditionals,
	)
}

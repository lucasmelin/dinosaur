package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Header implements a DNS message header as defined in [RFC 1035 section 4.1.1].
//
// [RFC 1035 section 4.1.1]: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.1
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
	//  - Z - a three bit field that is reserved for future use. Must be zero in all queries and responses.
	//  - RCODE - a four bit field that is set as part of responses. The values are no error condition (0),
	//    a format error (1) meaning the name server was unable to interpret the query, a server failure (2) meaning
	//    the name server was unable to process the query due to a problem with the name server, a name error (3) meaning
	//    that the domain name referenced in the query does not exist, not implemented (4) meaning the name server does
	//    not support the requested kind of query, refused (5) meaning the name server refuses to perform the specified
	//    operation for policy reasons, reserved for future use (6-15).
	Flags uint16
	// NumQuestions represents QDCOUNT, an unsigned 16-bit integer specifying the
	// number of entries in the question section.
	NumQuestions uint16
	// NumAnswers represents ANCOUNT, an unsigned 16-bit integer specifying the
	// number of resource records in the answer section.
	NumAnswers uint16
	// NumAuthorities represents NSCOUNT, an unsigned 16-bit integer specifying the
	// number of name server resource records in the authority records section.
	NumAuthorities uint16
	// NumAdditionals represents ARCOUNT, an unsigned 16-bit integer specifying the
	// number of resource records in the additional records section.
	NumAdditionals uint16
}

// ParseHeader parses a given bytes reader into a Header.
func ParseHeader(reader *bytes.Reader) Header {
	h := Header{}
	binary.Read(reader, binary.BigEndian, &h.ID)
	binary.Read(reader, binary.BigEndian, &h.Flags)
	binary.Read(reader, binary.BigEndian, &h.NumQuestions)
	binary.Read(reader, binary.BigEndian, &h.NumAnswers)
	binary.Read(reader, binary.BigEndian, &h.NumAuthorities)
	binary.Read(reader, binary.BigEndian, &h.NumAdditionals)

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

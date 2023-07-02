package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Question struct {
	// Name represents the query domain name.
	Name []byte
	// Type represents the query type.
	Type uint16
	// Class represents the query class.
	Class uint16
}

// ParseQuestion parses a given bytes reader into a Question.
func ParseQuestion(reader *bytes.Reader) Question {
	q := Question{}
	q.Name = DecodeName(reader)

	binary.Read(reader, binary.BigEndian, &q.Type)
	binary.Read(reader, binary.BigEndian, &q.Class)

	return q
}

// ToBytes encodes a Question as bytes.
// In network messages, integers are always encoded in big endian.
func (q *Question) ToBytes() []byte {
	b := make([]byte, 0)
	b = append(b, EncodeName(string(q.Name))...)
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, q.Type)
	_ = binary.Write(buf, binary.BigEndian, q.Class)
	b = append(b, buf.Bytes()...)
	return b
}

// EncodeName encodes a domain name by splitting the domain
// name into parts (labels), with each part prepended with its length.
// The encoded name is then returned as a byte array.
//
// This encoding format is defined in [RFC 1035 section 3.3].
//
// [RFC 1035 section 3.3]: https://datatracker.ietf.org/doc/html/rfc1035#section-3.3
func EncodeName(domainName string) []byte {
	parts := strings.Split(domainName, ".")
	buf := new(bytes.Buffer)
	for _, part := range parts {
		// Prepend the length of the part
		_ = binary.Write(buf, binary.BigEndian, uint8(len(part)))
		// Add the encoded part
		buf.Write([]byte(part))
	}
	// Add a zero byte to the end
	_ = binary.Write(buf, binary.BigEndian, uint8(0))
	return buf.Bytes()
}

func (q Question) String() string {
	return fmt.Sprintf(`Question{
    Name: %s,
    Type: %d,
    Class: %d
  }`, q.Name, q.Type, q.Class)
}

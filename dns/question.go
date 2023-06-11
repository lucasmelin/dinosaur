package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

type Question struct {
	Name  []byte
	Type  uint16
	Class uint16
}

func ParseQuestion(reader *bytes.Reader) Question {
	q := Question{}
	q.Name = DecodeName(reader)

	binary.Read(reader, binary.BigEndian, &q.Type)
	binary.Read(reader, binary.BigEndian, &q.Class)

	return q
}

// ToBytes encodes a Question as bytes.
// In network packets, integers are always encoded in big endian.
func (q *Question) ToBytes() []byte {
	b := make([]byte, 0)
	b = append(b, EncodeDNSName(string(q.Name))...)
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, q.Type)
	_ = binary.Write(buf, binary.BigEndian, q.Class)
	b = append(b, buf.Bytes()...)
	return b
}

func EncodeDNSName(domainName string) []byte {
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
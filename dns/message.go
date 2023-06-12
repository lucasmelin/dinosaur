package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

// Message represents a full DNS message.
// See: https://datatracker.ietf.org/doc/html/rfc1035#section-4
type Message struct {
	// header includes fields that specify which of the remaining fields are present,whether a message is a query
	// or a response, a standard query or some other opcode, etc.
	header    Header
	questions []Question
	// answers contains records that answer the question.
	answers []Record
	// authorities contains records that point towards an authoritative name server.
	authorities []Record
	// additionals contains records that relate to the query, but are not strictly answers to the question.
	additionals []Record
}

// ParseMessage parses a given byte array into a Message.
func ParseMessage(data []byte) Message {
	reader := bytes.NewReader(data)
	header := ParseHeader(reader)

	questions := []Question{}
	var i uint16
	for i = 0; i < header.NumQuestions; i++ {
		questions = append(questions, ParseQuestion(reader))
	}
	answers := []Record{}
	for i = 0; i < header.NumAnswers; i++ {
		answers = append(answers, ParseRecord(reader))
	}
	authorities := []Record{}
	for i = 0; i < header.NumAuthorities; i++ {
		authorities = append(authorities, ParseRecord(reader))
	}
	additionals := []Record{}
	for i = 0; i < header.NumAdditionals; i++ {
		additionals = append(additionals, ParseRecord(reader))
	}
	return Message{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

// DecodeName returns the first domain name found in the provided reader.
func DecodeName(reader *bytes.Reader) []byte {
	var parts = []string{}
	for length, _ := reader.ReadByte(); int(length) != 0; length, _ = reader.ReadByte() {

		// Check if the first 2 bits are 1s.
		// Any length starting with 11 means that the name is compressed.
		if length&0b1100_0000 == 0b1100_0000 {
			decompressed := DecodeCompressedName(length, reader)
			parts = append(parts, string(decompressed))
			// A compressed name is never followed by another label
			break
		} else {
			// Make a new byte array to read
			rdbuf := make([]byte, length)
			_, _ = reader.Read(rdbuf)
			parts = append(parts, string(rdbuf))
		}
	}
	return []byte(strings.Join(parts, "."))
}

// DecodeCompressedName extracts a domain name from a compressed message response.
// See: https://datatracker.ietf.org/doc/html/rfc1035#section-4.1.4
func DecodeCompressedName(length byte, reader *bytes.Reader) []byte {
	// Read a single byte
	b, _ := reader.ReadByte()
	// Take the bottom 6 bits of the length byte plus the next byte
	pointerBytes := [2]byte{length & 0b00111111, b}

	var pointer uint16
	binary.Read(bytes.NewBuffer(pointerBytes[:]), binary.BigEndian, &pointer)

	currentPosition, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return []byte("")
	}

	_, err = reader.Seek(int64(pointer), io.SeekStart)
	if err != nil {
		return []byte("")
	}
	result := DecodeName(reader)

	// Restore the current position in reader
	_, err = reader.Seek(currentPosition, io.SeekStart)
	if err != nil {
		return []byte("")
	}

	return result
}

// String formats a Message for printing.
func (p Message) String() string {
	return fmt.Sprintf(
		`Message{
  header: %s,
  questions: %+v,
  answers: %+v,
  authorities: %+v,
  additionals: %+v
}`,
		p.header,
		p.questions,
		p.answers,
		p.authorities,
		p.additionals,
	)
}

// GetAnswer returns the Data field from the first A record answer field in the Message.
func GetAnswer(message Message) []byte {
	for _, answer := range message.answers {
		if answer.Type == TypeA {
			return answer.Data
		}
	}
	return nil
}

// GetNameserverIP returns the Data field from the first A record additional field in the Message.
func GetNameserverIP(message Message) []byte {
	for _, additional := range message.additionals {
		if additional.Type == TypeA {
			return additional.Data
		}
	}
	return nil
}

// GetNameserver returns the Data field from the first NS record authority field in the Message.
func GetNameserver(message Message) string {
	for _, authority := range message.authorities {
		if authority.Type == TypeNS {
			return string(authority.Data)
		}
	}
	return ""
}

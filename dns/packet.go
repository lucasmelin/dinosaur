package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
)

type Packet struct {
	header      Header
	questions   []Question
	answers     []Record
	authorities []Record
	additionals []Record
}

func ParseDNSPacket(data []byte) Packet {
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
	return Packet{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

func DecodeName(reader *bytes.Reader) []byte {
	var parts = []string{}
	for length, _ := reader.ReadByte(); int(length) != 0; length, _ = reader.ReadByte() {

		// Check if the first 2 bits are 1s
		if length&0b1100_0000 == 0b1100_0000 {
			decompressed := DecodeCompressedName(length, reader)
			parts = append(parts, string(decompressed))
			// A compressed Name is never followed by another label
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

func (p Packet) String() string {
	return fmt.Sprintf(
		`Packet{
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

func GetAnswer(packet Packet) []byte {
	for _, answer := range packet.answers {
		if answer.Type == TypeA {
			return answer.Data
		}
	}
	return nil
}

func GetNameserverIP(packet Packet) []byte {
	for _, additional := range packet.additionals {
		if additional.Type == TypeA {
			return additional.Data
		}
	}
	return nil
}

func GetNameserver(packet Packet) string {
	for _, authority := range packet.authorities {
		if authority.Type == TypeNS {
			return string(authority.Data)
		}
	}
	return ""
}

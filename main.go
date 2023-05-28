package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	ClassIn = 1
	TypeA   = 1
)

type DNSHeader struct {
	id             uint16
	flags          uint16
	numQuestions   uint16
	numAnswers     uint16
	numAuthorities uint16
	numAdditionals uint16
}

type DNSQuestion struct {
	name  string
	Type  uint16
	Class uint16
}

type DNSRecord struct {
	// name represents the domain name.
	name []byte
	// A record type such as A, AAAA, MX, NS, TXT, etc.
	Type  uint16
	Class uint16
	// ttl represents the time-to-live of the cache entry for the query.
	ttl int32
	// data contains the records content, such as the IP address.
	data []byte
}

type DNSPacket struct {
	header      DNSHeader
	questions   []DNSQuestion
	answers     []DNSRecord
	authorities []DNSRecord
	additionals []DNSRecord
}

// toBytes encodes a DNSHeader as bytes.
// In network packets, integers are always encoded in big endian.
func (h *DNSHeader) toBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, h.id)
	_ = binary.Write(buf, binary.BigEndian, h.flags)
	_ = binary.Write(buf, binary.BigEndian, h.numQuestions)
	_ = binary.Write(buf, binary.BigEndian, h.numAnswers)
	_ = binary.Write(buf, binary.BigEndian, h.numAuthorities)
	_ = binary.Write(buf, binary.BigEndian, h.numAdditionals)
	return buf.Bytes()
}

// toBytes encodes a DNSQuestion as bytes.
// In network packets, integers are always encoded in big endian.
func (q *DNSQuestion) toBytes() []byte {
	b := make([]byte, 0)
	b = append(b, encodeDNSName(q.name)...)
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, q.Type)
	_ = binary.Write(buf, binary.BigEndian, q.Class)
	b = append(b, buf.Bytes()...)
	return b
}

func encodeDNSName(domainName string) []byte {
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

func buildQuery(queryID int, domainName string, recordType uint16) []byte {
	RecursionDesired := 1 << 8
	header := DNSHeader{
		id:             uint16(queryID),
		flags:          uint16(RecursionDesired),
		numQuestions:   1,
		numAnswers:     0,
		numAuthorities: 0,
		numAdditionals: 0,
	}
	question := DNSQuestion{
		name:  domainName,
		Type:  recordType,
		Class: ClassIn,
	}
	return append(header.toBytes(), question.toBytes()...)
}

func randomID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(65535)
}

func sendQuery(domain string) {
	query := buildQuery(randomID(), domain, TypeA)
	con, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	_, _ = con.Write(query)
	response := make([]byte, 1024)
	_, err = con.Read(response)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("./response", response, 0644)
	fmt.Println(parseDNSPacket(response))
}

func parseHeader(header *bytes.Reader) DNSHeader {
	h := DNSHeader{}
	binary.Read(header, binary.BigEndian, &h.id)
	binary.Read(header, binary.BigEndian, &h.flags)
	binary.Read(header, binary.BigEndian, &h.numQuestions)
	binary.Read(header, binary.BigEndian, &h.numAnswers)
	binary.Read(header, binary.BigEndian, &h.numAuthorities)
	binary.Read(header, binary.BigEndian, &h.numAdditionals)

	return h
}

func (h DNSHeader) String() string {
	return fmt.Sprintf(
		`DNSHeader{
    id: %d,
    flags: %d,
    numQuestions: %d,
    numAnswers: %d,
    numAuthorities: %d,
    numAdditionals: %d
  }`,
		h.id,
		h.flags,
		h.numQuestions,
		h.numAnswers,
		h.numAuthorities,
		h.numAdditionals,
	)
}

func decodeName(reader *bytes.Reader) string {
	var parts = []string{}
	for length, _ := reader.ReadByte(); int(length) != 0; length, _ = reader.ReadByte() {

		// Check if the first 2 bits are 1s
		if length&0b1100_0000 == 0b1100_0000 {
			decompressed := decodeCompressedName(length, reader)
			parts = append(parts, decompressed)
			// A compressed name is never followed by another label
			break
		} else {
			// Make a new byte array to read
			rdbuf := make([]byte, length)
			_, _ = reader.Read(rdbuf)
			parts = append(parts, string(rdbuf))
		}
	}
	return strings.Join(parts, ".")
}

func decodeCompressedName(length byte, reader *bytes.Reader) string {
	// Read a single byte
	b, _ := reader.ReadByte()
	// Take the bottom 6 bits of the length byte plus the next byte
	pointerBytes := [2]byte{length & 0b00111111, b}

	var pointer uint16
	binary.Read(bytes.NewBuffer(pointerBytes[:]), binary.BigEndian, &pointer)

	currentPosition, err := reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return ""
	}

	_, err = reader.Seek(int64(pointer), io.SeekStart)
	if err != nil {
		return ""
	}
	result := decodeName(reader)

	// Restore the current position in reader
	_, err = reader.Seek(currentPosition, io.SeekStart)
	if err != nil {
		return ""
	}

	return result
}

func parseQuestion(reader *bytes.Reader) DNSQuestion {
	q := DNSQuestion{}
	q.name = decodeName(reader)

	binary.Read(reader, binary.BigEndian, &q.Type)
	binary.Read(reader, binary.BigEndian, &q.Class)

	return q
}

func (q DNSQuestion) String() string {
	return fmt.Sprintf(`DNSQuestion{
    name: %s,
    Type: %d,
    Class: %d
  }`, q.name, q.Type, q.Class)
}

func parseRecord(br *bytes.Reader) DNSRecord {
	record := DNSRecord{}
	name := decodeName(br)

	record.name = []byte(name)

	binary.Read(br, binary.BigEndian, &record.Type)
	binary.Read(br, binary.BigEndian, &record.Class)
	binary.Read(br, binary.BigEndian, &record.ttl)

	var dataLength uint16
	binary.Read(br, binary.BigEndian, &dataLength)
	record.data = make([]byte, dataLength)
	binary.Read(br, binary.BigEndian, &record.data)

	return record
}

func (r DNSRecord) String() string {
	return fmt.Sprintf(`DNSRecord{
    name: %s,
    Type: %d,
    Class: %d,
    TTL: %d,
    Data: %s
  }`, r.name, r.Type, r.Class, r.ttl, IPString(r.data))
}

func parseDNSPacket(data []byte) DNSPacket {
	reader := bytes.NewReader(data)
	header := parseHeader(reader)

	questions := []DNSQuestion{}
	var i uint16
	for i = 0; i < header.numQuestions; i++ {
		questions = append(questions, parseQuestion(reader))
	}
	answers := []DNSRecord{}
	for i = 0; i < header.numAnswers; i++ {
		answers = append(answers, parseRecord(reader))
	}
	authorities := []DNSRecord{}
	for i = 0; i < header.numAuthorities; i++ {
		authorities = append(authorities, parseRecord(reader))
	}
	additionals := []DNSRecord{}
	for i = 0; i < header.numAdditionals; i++ {
		additionals = append(additionals, parseRecord(reader))
	}
	return DNSPacket{
		header:      header,
		questions:   questions,
		answers:     answers,
		authorities: authorities,
		additionals: additionals,
	}
}

func (p DNSPacket) String() string {
	return fmt.Sprintf(
		`DNSPacket{
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

func IPString(data []byte) string {
	if len(data) < 4 {
		return "?.?.?.?"
	}
	return fmt.Sprintf("%v.%v.%v.%v", data[0], data[1], data[2], data[3])
}

func main() {
	sendQuery("www.recurse.com")
}

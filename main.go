package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"strings"
)

const (
	ClassIn = 1
	TypeA   = 1
)

type DNSHeader struct {
	id             int
	flags          int
	numQuestions   int
	numAnswers     int
	numAuthorities int
	numAdditionals int
}

type DNSQuestion struct {
	name  string
	Type  int
	Class int
}

// toBytes encodes a DNSHeader as bytes.
// In network packets, integers are always encoded in big endian.
func (h *DNSHeader) toBytes() []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, uint16(h.id))
	_ = binary.Write(buf, binary.BigEndian, uint16(h.flags))
	_ = binary.Write(buf, binary.BigEndian, uint16(h.numQuestions))
	_ = binary.Write(buf, binary.BigEndian, uint16(h.numAnswers))
	_ = binary.Write(buf, binary.BigEndian, uint16(h.numAuthorities))
	_ = binary.Write(buf, binary.BigEndian, uint16(h.numAdditionals))
	return buf.Bytes()
}

// toBytes encodes a DNSQuestion as bytes.
// In network packets, integers are always encoded in big endian.
func (q *DNSQuestion) toBytes() []byte {
	b := make([]byte, 0)
	b = append(b, encodeDNSName(q.name)...)
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, uint16(q.Type))
	_ = binary.Write(buf, binary.BigEndian, uint16(q.Class))
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

func buildQuery(queryID int, domainName string, recordType int) []byte {
	RecursionDesired := 1 << 8
	header := DNSHeader{
		id:             queryID,
		flags:          RecursionDesired,
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
	return rand.Intn(65535)
}

func sendQuery(domain string) {
	query := buildQuery(randomID(), domain, TypeA)
	con, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	con.Write(query)
	response := make([]byte, 1024)
	_, err = con.Read(response)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	sendQuery("www.example.com")
}

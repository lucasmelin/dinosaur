package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
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

type DNSRecord struct {
	// name represents the domain name.
	name []byte
	// A record type such as A, AAAA, MX, NS, TXT, etc.
	Type  int
	Class int
	// ttl represents the time-to-live of the cache entry for the query.
	ttl int
	// data contains the records content, such as the IP address.
	data []byte
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

	_, _ = con.Write(query)
	response := make([]byte, 1024)
	_, err = con.Read(response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q", response)
	fmt.Println(parseHeader(response))
}

func parseHeader(header []byte) DNSHeader {
	fields := make([]int, 0)
	// Each of the 6 fields is a 2-byte integer.
	for i := 0; i < 12; i += 2 {
		val := binary.BigEndian.Uint16(header[i : i+2])
		fields = append(fields, int(val))
	}
	return DNSHeader{
		id:             fields[0],
		flags:          fields[1],
		numQuestions:   fields[2],
		numAnswers:     fields[3],
		numAuthorities: fields[4],
		numAdditionals: fields[5],
	}
}

func (h DNSHeader) String() string {
	return fmt.Sprintf(
		"DNSHeader{id: %d, flags: %d, numQuestions: %d, numAnswers: %d, numAuthorities: %d, numAdditionals: %d}",
		h.id,
		h.flags,
		h.numQuestions,
		h.numAnswers,
		h.numAuthorities,
		h.numAdditionals,
	)
}

func main() {
	sendQuery("www.example.com")
}

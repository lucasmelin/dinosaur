package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	ClassIn = 1
	// IP Address for a.root-servers.net
	rootNameserver = "198.41.0.4"
)

const (
	TypeA = iota + 1
	TypeNS
	TypeMD
	TypeMF
	TypeCNAME
	TypeSOA
	TypeMB
	TypeMG
	TypeMR
	TypeNULL
	TypeWKS
	TypePTR
	TypeHINFO
	TypeMINFO
	TypeMX
	TypeTXT
)

type RecordType struct {
	Name    string
	Value   uint16
	Meaning string
}

// RecordTypes represents all possible resource record TYPE field values.
// See: https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
var RecordTypes = map[string]RecordType{
	"A": {
		Name:    "A",
		Value:   TypeA,
		Meaning: "a host address",
	},
	"NS": {
		Name:    "NS",
		Value:   TypeNS,
		Meaning: "an authoritative name server",
	},
	"MD": {
		Name:    "MD",
		Value:   TypeMD,
		Meaning: "a mail destination (Obsolete - use MX)",
	},
	"MF": {
		Name:    "MF",
		Value:   TypeMF,
		Meaning: "a mail forwarder (Obsolete - use MX)",
	},
	"CNAME": {
		Name:    "CNAME",
		Value:   TypeCNAME,
		Meaning: "the canonical name for an alias",
	},
	"SOA": {
		Name:    "SOA",
		Value:   TypeSOA,
		Meaning: "marks the start of a zone of authority",
	},
	"MB": {
		Name:    "MB",
		Value:   TypeMB,
		Meaning: "a mailbox domain name (EXPERIMENTAL)",
	},
	"MG": {
		Name:    "MG",
		Value:   TypeMG,
		Meaning: "a mail group member (EXPERIMENTAL)",
	},
	"MR": {
		Name:    "MR",
		Value:   TypeMR,
		Meaning: "a mail rename domain name (EXPERIMENTAL)",
	},
	"NULL": {
		Name:    "NULL",
		Value:   TypeNULL,
		Meaning: "a null RR (EXPERIMENTAL)",
	},
	"WKS": {
		Name:    "WKS",
		Value:   TypeWKS,
		Meaning: "a well known service description",
	},
	"PTR": {
		Name:    "PTR",
		Value:   TypePTR,
		Meaning: "a domain name pointer",
	},
	"HINFO": {
		Name:    "HINFO",
		Value:   TypeHINFO,
		Meaning: "host information",
	},
	"MINFO": {
		Name:    "MINFO",
		Value:   TypeMINFO,
		Meaning: "mailbox or mail list information",
	},
	"MX": {
		Name:    "MX",
		Value:   TypeMX,
		Meaning: "mail exchange",
	},
	"TXT": {
		Name:    "TXT",
		Value:   TypeTXT,
		Meaning: "text strings",
	},
}

// Record represents a DNS resource record.
// See: https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.1
type Record struct {
	// Name represents the domain Name.
	Name []byte
	// A record type such as A, AAAA, MX, NS, TXT, etc.
	Type  uint16
	Class uint16
	// TTL represents the time-to-live of the cache entry for the query.
	TTL int32
	// Data contains the records content, such as the IP address.
	Data []byte
}

// ParseRecord parses a given bytes reader into a record.
func ParseRecord(reader *bytes.Reader) Record {
	record := Record{}

	record.Name = DecodeName(reader)

	binary.Read(reader, binary.BigEndian, &record.Type)
	binary.Read(reader, binary.BigEndian, &record.Class)
	binary.Read(reader, binary.BigEndian, &record.TTL)

	var dataLength uint16
	binary.Read(reader, binary.BigEndian, &dataLength)
	data := make([]byte, dataLength)
	// Save the offset before reading the data field.
	beforeData, _ := reader.Seek(0, io.SeekCurrent)
	binary.Read(reader, binary.BigEndian, &data)

	switch record.Type {
	case TypeNS:
		record.Data = DecodeName(bytes.NewReader(data))
	case TypeA:
		record.Data = []byte(IPString(data))
	case TypeCNAME:
		// Seek back in the reader to before the name, so that
		// DecodeName can decompress the name by referring to bytes
		// anywhere in the response.
		reader.Seek(beforeData, io.SeekStart)
		record.Data = DecodeName(reader)
	default:
		record.Data = data
	}

	return record
}

// IPString converts a byte array into a dotted IP format.
func IPString(data []byte) string {
	if len(data) < 4 {
		return "?.?.?.?"
	}
	return fmt.Sprintf("%v.%v.%v.%v", data[0], data[1], data[2], data[3])
}

func (r Record) String() string {
	data := r.Data
	if r.Type == TypeA {
		data = []byte(IPString(r.Data))
	}

	return fmt.Sprintf(`Record{
    Name: %s,
    Type: %d,
    Class: %d,
    TTL: %d,
    Data: %s
  }`, r.Name, r.Type, r.Class, r.TTL, data)
}

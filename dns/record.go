package dns

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	ClassIn        = 1
	TypeA          = 1
	TypeNS         = 2
	TypeTxt        = 16
	rootNameserver = ""
)

type RecordType struct {
	Name    string
	Value   int
	Meaning string
}

// RecordTypes represents all possible resource record TYPE field values.
// See: https://datatracker.ietf.org/doc/html/rfc1035#section-3.2.2
var RecordTypes = []RecordType{
	{
		Name:    "A",
		Value:   1,
		Meaning: "a host address",
	},
	{
		Name:    "NS",
		Value:   2,
		Meaning: "an authoritative name server",
	},
}

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

func ParseRecord(br *bytes.Reader) Record {
	record := Record{}

	record.Name = DecodeName(br)

	binary.Read(br, binary.BigEndian, &record.Type)
	binary.Read(br, binary.BigEndian, &record.Class)
	binary.Read(br, binary.BigEndian, &record.TTL)

	var dataLength uint16
	binary.Read(br, binary.BigEndian, &dataLength)
	data := make([]byte, dataLength)
	binary.Read(br, binary.BigEndian, &data)

	if record.Type == TypeNS {
		record.Data = DecodeName(bytes.NewReader(data))
	} else if record.Type == TypeA {
		record.Data = []byte(IPString(data))
	} else {
		record.Data = data
	}

	return record
}

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

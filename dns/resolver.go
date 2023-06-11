package dns

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	RecursionDesired = 1 << 8
	RecursionOff     = 0
)

func BuildQuery(queryID int, domainName string, recordType uint16) []byte {
	header := Header{
		Id:             uint16(queryID),
		Flags:          RecursionOff,
		NumQuestions:   1,
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}
	question := Question{
		Name:  []byte(domainName),
		Type:  recordType,
		Class: ClassIn,
	}
	return append(header.ToBytes(), question.ToBytes()...)
}

func RandomID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(65535)
}

func SendQuery(ipAddress string, domain string, recordType uint16) Packet {
	query := BuildQuery(RandomID(), domain, recordType)
	con, err := net.Dial("udp", fmt.Sprintf("%s:53", ipAddress))
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
	return ParseDNSPacket(response)
}

func Resolve(domainName string, recordType uint16) []byte {
	nameserver := "198.41.0.4"
	for true {
		fmt.Fprintf(os.Stdout, "Querying %s for %s\n", nameserver, domainName)
		response := SendQuery(nameserver, domainName, recordType)
		if ip := GetAnswer(response); ip != nil {
			return ip
		} else if nsIP := GetNameserverIP(response); nsIP != nil {
			nameserver = string(nsIP)
		} else if nsDomain := GetNameserver(response); nsDomain != "" {
			nameserver = string(Resolve(nsDomain, TypeA))
		} else {
			panic("something went wrong")
		}
	}
	return nil
}
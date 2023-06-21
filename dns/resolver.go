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

func BuildQuery(queryID int, domainName string, recordType string) []byte {
	recType := RecordTypes[recordType]
	header := Header{
		ID:             uint16(queryID),
		Flags:          RecursionOff,
		NumQuestions:   1,
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}
	question := Question{
		Name:  []byte(domainName),
		Type:  recType.Value,
		Class: ClassIn,
	}
	return append(header.ToBytes(), question.ToBytes()...)
}

// RandomID returns a random 16-bit integer.
func RandomID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(65535)
}

// SendQuery sends a query to a given DNS resolver, and returns the message from the resolver.
func SendQuery(ipAddress string, domain string, recordType string) Message {
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
	return ParseMessage(response)
}

// Resolve recursively queries nameservers to find the IP address for a given domain name.
func Resolve(domainName string, recordType string) []byte {
	nameserver := rootNameserver
	for true {
		fmt.Printf("Querying %s for %s\n", nameserver, domainName)
		response := SendQuery(nameserver, domainName, recordType)
		if ip := GetAnswer(response); ip != nil {
			return ip
		} else if nsIP := GetNameserverIP(response); nsIP != nil {
			nameserver = string(nsIP)
		} else if nsDomain := GetNameserver(response); nsDomain != "" {
			nameserver = string(Resolve(nsDomain, "A"))
		} else {
			panic("something went wrong")
		}
	}
	return nil
}

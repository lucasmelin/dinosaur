package dns

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/lucasmelin/dinosaur/dino"
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
	con, err := net.Dial("udp", fmt.Sprintf("%s:53", ipAddress))
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	query := BuildQuery(RandomID(), domain, recordType)
	_, _ = con.Write(query)
	response := make([]byte, 1024)
	_, err = con.Read(response)
	if err != nil {
		log.Fatal(err)
	}
	// _ = os.WriteFile("./response", response, 0644)
	return ParseMessage(response)
}

// AsciiResolve recursively queries nameservers to find the IP address for a given domain name.
// It also prints ascii art of each resolution step.
func AsciiResolve(domainName string, recordType string, impatient bool) []byte {
	nameserver := rootNameserver
	for {
		dino.NewDino().SayRight(fmt.Sprintf("Hey %s what's the address for %s?", nameserver, domainName))
		if !impatient {
			waitForKeypress()
		}
		response := SendQuery(nameserver, domainName, recordType)
		if ip := GetAnswer(response); ip != nil {
			dino.NewServer(nameserver).SayLeft(fmt.Sprintf("The IP address is %s", ip))
			if !impatient {
				waitForKeypress()
			}
			return ip
		} else if alias := GetAlias(response); alias != "" {
			dino.NewServer(nameserver).SayLeft(fmt.Sprintf("That's an alias for %s", alias))
			if !impatient {
				waitForKeypress()
			}
			return AsciiResolve(alias, "A", impatient)
		} else if nsIP := GetNameserverIP(response); nsIP != nil {
			dino.NewServer(nameserver).SayLeft(fmt.Sprintf("I don't know, you should ask %s", nsIP))
			if !impatient {
				waitForKeypress()
			}
			nameserver = string(nsIP)
		} else if nsDomain := GetNameserver(response); nsDomain != "" {
			dino.NewServer(nameserver).SayLeft(fmt.Sprintf("I don't know, you should ask %s", nsDomain))
			if !impatient {
				waitForKeypress()
			}
			dino.NewDino().SayLeft(fmt.Sprintf("I wonder how I can reach %s", nsDomain))
			if !impatient {
				waitForKeypress()
			}
			nameserver = string(AsciiResolve(nsDomain, "A", impatient))
		} else {
			panic("something went wrong")
		}
	}
	return nil
}

func waitForKeypress() {
	fmt.Println("\nPress any key to continue the journey...")
	_, _ = fmt.Scanln()
}

// Resolve recursively queries nameservers to find the IP address for a given domain name.
func Resolve(domainName string, recordType string) []byte {
	nameserver := rootNameserver
	for {
		response := SendQuery(nameserver, domainName, recordType)
		if ip := GetAnswer(response); ip != nil {
			return ip
		} else if alias := GetAlias(response); alias != "" {
			return Resolve(alias, "A")
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

package main

import (
	"fmt"

	"github.com/lucasmelin/dinosaur/dns"
)

func main() {
	fmt.Println(dns.Resolve("twitter.com", dns.TypeA))
	// fmt.Println(sendQuery("198.41.0.4", "www.recurse.com", TypeA))
}

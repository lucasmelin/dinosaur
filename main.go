package main

import (
	"fmt"

	"github.com/lucasmelin/dinosaur/dns"
)

func main() {
	fmt.Println(string(dns.Resolve("www.facebook.com", "A")))
	// fmt.Println(sendQuery("198.41.0.4", "www.recurse.com", TypeA))
}

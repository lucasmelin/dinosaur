package main

import (
	"fmt"

	"github.com/lucasmelin/dinosaur/dino"
	"github.com/lucasmelin/dinosaur/dns"
)

func main() {
	url := "www.lucasmelin.com"
	d := dino.NewDino()
	d.SayLeft(fmt.Sprintf("I wonder how I can reach %s", url))
	ipAddress := string(dns.Resolve(url, "A"))
	d.SayLeft(fmt.Sprintf("Great, now I know I can reach %s at %s", url, ipAddress))
}

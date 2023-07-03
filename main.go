package main

import (
	"flag"
	"fmt"

	"github.com/lucasmelin/dinosaur/dino"
	"github.com/lucasmelin/dinosaur/dns"
)

func main() {
	var meteor = flag.Bool("meteor", false, "disable the dino ascii art (default false)")
	var impatient = flag.Bool("impatient", false, "do not wait for an input between each frame (default false)")
	flag.Parse()
	args := flag.Args()
	url := "www.lucasmelin.com"
	if len(args) != 0 {
		url = args[0]
	}
	if !*meteor {
		d := dino.NewDino()
		d.SayLeft(fmt.Sprintf("I wonder how I can reach %s", url))
		fmt.Println("\nPress any key to continue the journey...")
		_, _ = fmt.Scanln()
		ipAddress := string(dns.AsciiResolve(url, "A", *impatient))
		d.SayLeft(fmt.Sprintf("Great, now I know I can reach %s at %s", url, ipAddress))
	} else {
		ipAddress := string(dns.Resolve(url, "A"))
		fmt.Println(ipAddress)
	}
}

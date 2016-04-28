package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	domain := os.Args[1]
	names, err := net.LookupIP(domain)
	if err != nil {
		fail(err.Error())
	}
	fmt.Printf("%+v", names)
}

func fail(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
	os.Exit(1)
}

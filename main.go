package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func gopos(ip, startPort, endPort string) {
	startPortInt, _ := strconv.Atoi(startPort)
	endPortInt, _ := strconv.Atoi(endPort)

	// Scan ports
	for port := startPortInt; port < endPortInt; port += 1 {
		ip := fmt.Sprintf("%s:%d", ip, port)
		_, err := net.Dial("tcp", ip)
		if err == nil {
			fmt.Printf("Port %d open\n", port)
		}
	}
}

func main() {
	ip := os.Args[1]
	startPort := os.Args[2]
	endPort := os.Args[3]
	gopos(ip, startPort, endPort)
}

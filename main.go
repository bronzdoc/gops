package main

import (
	"fmt"
	"github.com/bronzdoc/gops/lib/util"
	"net"
	"os"
	"strconv"
)

func gops(ip, startPort, endPort, portType string) {
	startPortInt, _ := strconv.Atoi(startPort)
	endPortInt, _ := strconv.Atoi(endPort)

	// Scan ports
	for port := startPortInt; port < endPortInt; port += 1 {
		ip := fmt.Sprintf("%s:%d", ip, port)
		_, err := net.Dial(portType, ip)

		if err == nil {
			if val, ok := util.CommonPorts[port]; ok {
				fmt.Printf("%s/%d open -- %s \n", portType, port, val)
			} else {
				fmt.Printf("%s/%d open -- N\\A\n", portType, port)
			}
		}
	}
}

func main() {
	ip := os.Args[1]
	startPort := os.Args[2]
	endPort := os.Args[3]
	gops(ip, startPort, endPort, "tcp")
}

package main

import (
	"fmt"
	"github.com/bronzdoc/gops/lib/util"
	"github.com/gosuri/uitable"
	"net"
	"os"
)

func gops(ip, protocol string) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("PORT", "PROTOCOL", "DESCRIPTION")

	// Scan ports
	for port := 0; port <= 65535; port += 1 {
		ip := fmt.Sprintf("%s:%d", ip, port)
		_, err := net.Dial(protocol, ip)
		if err == nil {
			if val, ok := util.CommonPorts[port]; ok {
				table.AddRow(port, protocol, val)
			} else {
				table.AddRow(port, protocol, "(unknown)")
			}
		}
	}
	fmt.Println(table)
}

func main() {
	ip := os.Args[1]
	gops(ip, "tcp")
}

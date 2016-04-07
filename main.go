package main

import (
	"fmt"
	"github.com/CrowdSurge/banner"
	"github.com/bronzdoc/gops/lib/util"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"net"
	"os"
	"time"
)

const PORTS = 65535

func gops(ip, protocol string) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("PORT", "PROTOCOL", "DESCRIPTION")

	status := uilive.New()
	status.Start()

	// Scan ports
	for port := 0; port <= PORTS; port++ {
		ip := fmt.Sprintf("%s:%d", ip, port)
		_, err := net.Dial(protocol, ip)
		if err == nil {
			if val, ok := util.CommonPorts[port]; ok {
				table.AddRow(port, protocol, val)
			} else {
				table.AddRow(port, protocol, "(unknown)")
			}
			fmt.Fprintf(status, "Scanning...(%d%%)\n", int((float32(port)/PORTS)*100))
			time.Sleep(time.Millisecond * 1)
		}
	}

	fmt.Fprintf(status, "Finished: Scanning (100%%)\n")
	status.Stop()
	fmt.Println(table)
}

func main() {
	banner.Print("gops")
	fmt.Println("")
	ip := os.Args[1]
	gops(ip, "tcp")
}

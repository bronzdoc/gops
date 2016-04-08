package main

import (
	"flag"
	"fmt"
	"github.com/bronzdoc/gops/lib/util"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"net"
)

func getProtocol(tcp, udp *bool) string {
	var protocol string
	if *tcp {
		protocol = "tcp"
	} else if *udp {
		protocol = "udp"
	} else {
		protocol = "all"
	}
	return protocol
}

func scanTCP(host string, port int, table *uitable.Table) {
	_, err := net.Dial("tcp", host)
	if err == nil {
		if val, ok := util.CommonPorts[port]; ok {
			table.AddRow(port, "tcp", val)
		} else {
			table.AddRow(port, "tcp", "(unknown)")
		}
	}
}

func scanUDP(host string, port int, table *uitable.Table) {
	_, err := net.Dial("udp", host)
	if err == nil {
		if val, ok := util.CommonPorts[port]; ok {
			table.AddRow(port, "udp", val)
		} else {
			table.AddRow(port, "udp", "(unknown)")
		}
	}
}

func displayScanInfo(host string, port int, protocol string, table *uitable.Table) {
	if protocol == "tcp" {
		scanTCP(host, port, table)
	} else if protocol == "udp" {
		scanUDP(host, port, table)
	} else {
		scanTCP(host, port, table)
		scanUDP(host, port, table)
	}
}

func gops(options map[string]interface{}) {
	protocol := getProtocol(options["tcp"].(*bool), options["udp"].(*bool))
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("PORT", "PROTOCOL", "DESCRIPTION")

	status := uilive.New()
	status.Start()
	start := *options["start"].(*int)
	end := *options["end"].(*int)

	// Scan ports
	for port := start; port <= end; port++ {
		host := fmt.Sprintf("%s:%d", *options["host"].(*string), port)
		displayScanInfo(host, port, protocol, table)
		fmt.Fprintf(status, "gops scanning...(%d%%)\n", int((float32(port)/float32(end))*100))
		status.Flush()
	}

	fmt.Fprintf(status, "gops finished scanning (100%%)\n")
	status.Stop()
	fmt.Println(table)
}

func main() {
	options := map[string]interface{}{
		"host":  flag.String("host", "localhost", "host to scan"),
		"tcp":   flag.Bool("tcp", false, "Show only tcp ports open"),
		"udp":   flag.Bool("udp", false, "Show only udp ports open"),
		"start": flag.Int("start", 0, "Port to start the scan"),
		"end":   flag.Int("end", 65535, "Port to end the scan"),
	}
	flag.Parse()
	gops(options)
}

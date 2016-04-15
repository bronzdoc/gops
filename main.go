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

	if *tcp && *udp {
		protocol = "all"
	} else if *tcp {
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
			table.AddRow(port, "tcp", "(?)")
		}
	}
}

func scanUDP(host string, port int, table *uitable.Table) {
	serverAddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Write 3 times to the udp socket and check
	// if there's any kind of error
	error_count := 0
	for i := 0; i <= 3; i++ {
		buf := []byte("0")
		_, err := conn.Write(buf)
		if err != nil {
			error_count++
		}
	}

	if error_count <= 0 {
		if val, ok := util.CommonPorts[port]; ok {
			table.AddRow(port, "udp", val)
		} else {
			table.AddRow(port, "udp", "(?)")
		}
	}

	defer conn.Close()
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
	_ = status
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
		"host":  flag.String("host", "127.0.0.1", "host to scan"),
		"tcp":   flag.Bool("tcp", false, "Show only tcp ports open"),
		"udp":   flag.Bool("udp", false, "Show only udp ports open"),
		"start": flag.Int("start", 0, "Port to start the scan"),
		"end":   flag.Int("end", 65535, "Port to end the scan"),
	}
	flag.Parse()
	gops(options)
}

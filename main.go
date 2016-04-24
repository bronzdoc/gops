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

func scanTCP(host string, port int) int {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return -1
	}
	defer conn.Close()
	return port
}

func scanUDP(host string, port int) int {
	serverAddr, err := net.ResolveUDPAddr("udp", host)
	util.LogError(err)

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	util.LogError(err)

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	util.LogError(err)

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

	if error_count > 0 {
		return -1
	}
	defer conn.Close()
	return port
}

func displayScanInfo(host string, port int, protocol string, table *uitable.Table) {
	var udpPortScanned int
	var tcpPortScanned int
	var protocolDesc string

	if protocol == "tcp" {
		tcpPortScanned = scanTCP(host, port)
	} else if protocol == "udp" {
		udpPortScanned = scanUDP(host, port)
	} else {
		tcpPortScanned = scanTCP(host, port)
		udpPortScanned = scanUDP(host, port)
	}

	if tcpPortScanned != -1 || udpPortScanned != -1 {
		if tcpPortScanned == udpPortScanned {
			protocolDesc = "tcp/udp"
		} else if tcpPortScanned != -1 {
			protocolDesc = "tcp"
		} else if udpPortScanned != -1 {
			protocolDesc = "udp"
		}

		table.AddRow(port, protocolDesc, (func(port int) string {
			desc := "(?)"
			if val, ok := util.CommonPorts[port]; ok {
				desc = val
			}
			return desc
		}(port)))
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
		"host":  flag.String("host", "127.0.0.1", "host to scan"),
		"tcp":   flag.Bool("tcp", false, "Show only tcp ports open"),
		"udp":   flag.Bool("udp", false, "Show only udp ports open"),
		"start": flag.Int("start", 0, "Port to start the scan"),
		"end":   flag.Int("end", 65535, "Port to end the scan"),
	}
	flag.Parse()
	gops(options)
}

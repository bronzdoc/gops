package main

import (
	"flag"
	"fmt"
	"github.com/bronzdoc/gops/lib/util"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"net"
)

type ScanInfo struct {
	port     int
	protocol string
	desc     string
	empty    bool
}

type Job struct {
	host     string
	port     int
	end      int
	protocol string
}

func worker(statusChan chan int, jobs chan Job, resultsChan chan ScanInfo) {
	for job := range jobs {
		ScanInfo := getScannedInfo(job.host, job.port, job.protocol)
		if !ScanInfo.empty {
			resultsChan <- ScanInfo
		}
		statusChan <- 1
	}
}

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
	if err != nil {
		util.LogError(err)
		return -1
	}

	localAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		util.LogError(err)
		return -1
	}

	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		util.LogError(err)
		return -1
	}
	defer conn.Close()

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
	return port
}

func getScannedInfo(host string, port int, protocol string) ScanInfo {
	udpPortScanned := -1
	tcpPortScanned := -1
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
		info := ScanInfo{
			empty:    false,
			port:     port,
			protocol: protocolDesc,
			desc: func(port int) string { // Build protocol description
				desc := "(?)"
				if val, ok := util.CommonPorts[port]; ok {
					desc = val
				}
				return desc
			}(port),
		}
		return info
	}
	return ScanInfo{empty: true}
}

func gops() {

	resultsChan := make(chan ScanInfo, 10)
	jobs := make(chan Job, 10)
	statusChan := make(chan int)

	protocol := getProtocol(&tcp, &udp)

	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("PORT", "PROTOCOL", "DESCRIPTION")

	status := uilive.New()
	status.Start()

	display := uilive.New()
	display.Start()

	if port > 0 {
		start = port
		end = start + 1
	}

	portsToScann := end - start
	fmt.Printf("%d", portsToScann)
	scannedPorts := 0

	// Status handler
	go func() {
		for counter := range statusChan {
			scannedPorts += counter
			fmt.Fprintf(status, "gops scanning...(%d%%)\n", int((float32(scannedPorts)/float32(portsToScann))*100))
			status.Flush()
			if scannedPorts == portsToScann {
				close(resultsChan)
			}
		}
	}()

	// Workers
	for i := 0; i < 5; i++ {
		go worker(statusChan, jobs, resultsChan)
	}

	// Enqueue jobs
	for port := start; port <= end; port++ {
		host := fmt.Sprintf("%s:%d", host, port)
		jobs <- Job{host, port, end, protocol}
	}
	close(jobs)

	for Scannedinfo := range resultsChan {
		table.AddRow(Scannedinfo.port, Scannedinfo.protocol, Scannedinfo.desc)
	}

	fmt.Fprintf(status, "gops finished scanning (100%%)\n")
	status.Stop()

	fmt.Fprintf(display, "%s\n", table)
	display.Stop()
}

var (
	host  string
	tcp   bool
	udp   bool
	start int
	end   int
	port  int
)

func main() {
	flag.StringVar(&host, "host", "127.0.0.1", "Specify host")
	flag.BoolVar(&tcp, "tcp", false, "Show only tcp ports open")
	flag.BoolVar(&udp, "udp", false, "Show only udp ports open")
	flag.IntVar(&start, "start", 0, "Port to start the scan")
	flag.IntVar(&end, "end", 65535, "Port to end the scan")
	flag.IntVar(&port, "port", 0, "Check if port is open")

	flag.Parse()
	gops()
}

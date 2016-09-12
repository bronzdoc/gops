package main

import (
	"flag"
	"fmt"
	"github.com/bronzdoc/gops/lib/util"
	"github.com/gosuri/uilive"
	"github.com/gosuri/uitable"
	"github.com/tj/go-spin"
	"net"
	"sync"
	"time"
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

func worker(jobs <-chan Job, results chan<- ScanInfo, notifyFinish chan<- bool) {
	defer wg.Done()
	for job := range jobs {
		scanInfo := getScannedInfo(job.host, job.port, job.protocol)
		if !scanInfo.empty {
			results <- scanInfo
		}
	}
	notifyFinish <- true
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
		scanInfo := ScanInfo{
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
		return scanInfo
	}
	return ScanInfo{empty: true}
}

func gops() {
	results := make(chan ScanInfo, 10)
	jobs := make(chan Job, 10)
	notifyFinish := make(chan bool)

	MAXWORKERS := 10

	protocol := getProtocol(&tcp, &udp)

	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("PORT", "PROTOCOL", "DESCRIPTION")

	spinner := spin.New()

	display := uilive.New()
	display.Start()

	if port > 0 {
		start = port
		end = port
	}

	// loader handler
	go func() {
		for {
			fmt.Printf("\r  \033[36mgops %s\033[m ", spinner.Next())
			time.Sleep(60 * time.Millisecond)
		}
	}()

	// Workers
	for i := 0; i < MAXWORKERS; i++ {
		wg.Add(1)
		go worker(jobs, results, notifyFinish)
	}

	// Enqueue jobs
	for port := start; port <= end; port++ {
		host := fmt.Sprintf("%s:%d", host, port)
		jobs <- Job{host, port, end, protocol}
	}
	close(jobs)

	wg.Add(1)
	go func() {
		workersCount := MAXWORKERS
		defer wg.Done()
		for {
			select {
			case scannedInfo := <-results:
				table.AddRow(scannedInfo.port, scannedInfo.protocol, scannedInfo.desc)
			case <-notifyFinish:
				workersCount--
				if workersCount == 0 {
					return
				}
			}
		}
	}()

	wg.Wait()

	fmt.Fprintf(display, "\n%s\n", table)
	display.Stop()
}

var (
	host  string
	tcp   bool
	udp   bool
	start int
	end   int
	port  int
	wg    sync.WaitGroup
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

package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
)

func scanPort(host string, port int, wg *sync.WaitGroup, openPorts chan int) {
	defer wg.Done()
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err == nil {
		conn.Close()
		openPorts <- port
	}
}

func main() {
	var host string
	var minPort, maxPort int

	flag.StringVar(&host, "host", "", "Target host to scan")
	flag.IntVar(&minPort, "min", 0, "Minimum port number")
	flag.IntVar(&maxPort, "max", 0, "Maximum port number")
	flag.Parse()

	if host == "" || minPort == 0 || maxPort == 0 {
		fmt.Println("Please provide valid values for host, min, and max flags.")
		return
	}

	if minPort > maxPort {
		fmt.Println("Minimum port cannot be greater than maximum port.")
		return
	}

	openPorts := make(chan int, maxPort-minPort+1)
	var wg sync.WaitGroup

	for port := minPort; port <= maxPort; port++ {
		wg.Add(1)
		go scanPort(host, port, &wg, openPorts)
	}

	go func() {
		wg.Wait()
		close(openPorts)
	}()

	var result []int
	for openPort := range openPorts {
		result = append(result, openPort)
	}

	fmt.Printf("Open ports: %v\n", result)
}

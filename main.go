package main

import (
	"log"
	"os"

	"github.com/bongtavas/gonmap/lib/tcpscan"
)

func main() {
	if len(os.Args) != 3 {
		log.Printf("Usage: %s <host/ip> <port>\n", os.Args[0])
		os.Exit(-1)
	}
	log.Printf("gonmap %s %s", os.Args[1], os.Args[2])

	dstHostname := os.Args[1]
	dstPortStr := os.Args[2]

	tcpscan.SynScan(dstHostname, dstPortStr)
}

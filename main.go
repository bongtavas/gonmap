package main

import (
	"flag"
	"log"
	"os"

	"github.com/bongtavas/gonmap/lib/tcpscan"
	"github.com/bongtavas/gonmap/lib/udpscan"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: %s <host/ip>\n", os.Args[0])
		os.Exit(-1)
	}

	dstPorts := flag.String("p", "default", "Ports to scan")
	udpFlag := flag.Bool("udp", false, "Enable UDP Port Scanning")
	flag.Parse()

	dstHostname := flag.Args()[0]

	log.Println(*dstPorts)

	tcpscan.SynScan(dstHostname, *dstPorts)
	if *udpFlag {
		udpscan.UdpScan(dstHostname)
	}
}

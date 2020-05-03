package main

import (
	"flag"
	"log"
	"os"

	"github.com/bongtavas/gonmap/lib/netutil"
	"github.com/bongtavas/gonmap/lib/tcpscan"
	"github.com/bongtavas/gonmap/lib/udpscan"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("Usage: %s <host/ip>\n", os.Args[0])
		os.Exit(-1)
	}

	dstPorts := flag.String("p", "default", "Ports to scan")
	synScanFlag := flag.Bool("sS", false, "Enable TCP SYN Port Scanning")
	udpScanFlag := flag.Bool("sU", false, "Enable UDP Port Scanning")
	flag.Parse()

	dstHostname := flag.Args()[0]

	log.Println("Scanning ports " + *dstPorts)

	dstPortList := netutil.BuildPortList(*dstPorts)
	dstIPs, err := netutil.BuildIPList(dstHostname)

	if err != nil {
		log.Fatal("Build IP List failed, check your hostname argument")
	}

	if *synScanFlag {
		for _, dstIP := range dstIPs {
			tcpscan.SynScan(dstIP, dstPortList)
		}
	}

	if *udpScanFlag {
		for _, dstIP := range dstIPs {
			udpscan.UdpScan(dstIP)
		}
	}
}

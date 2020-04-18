package udpscan

import (
	"log"
	"strconv"

	"github.com/miekg/dns"
)

func UdpScan(dstHostname string) {
	dnsScan(dstHostname, 53)
}

// Since DNS is a UDP process, we need to send a valid DNS request to illicit a response from a service
// A response even if empty signifies presence of a DNS service
func dnsScan(dstHostname string, dstPort uint64) {
	target := "google.com"
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)

	log.Printf("[UDP Port %d] Sending DNS request to %s", dstPort, dstHostname)
	_, _, err := c.Exchange(&m, dstHostname+":"+strconv.FormatUint(dstPort, 10))
	if err != nil {
		log.Printf("[UDP Port %d] %v", dstPort, err)
		log.Printf("[UDP Port %d] is CLOSED", dstPort)
	} else {
		log.Printf("[UDP Port %d] DNS Response received", dstPort)
		log.Printf("[UDP Port %d] is OPEN", dstPort)
	}

	// No need to do anything with the DNS Response
}

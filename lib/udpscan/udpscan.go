package udpscan

import (
	"log"
	"os"
	"strconv"

	"github.com/miekg/dns"
	"github.com/olekukonko/tablewriter"
	"github.com/soniah/gosnmp"
)

func UdpScan(dstHostname string) {
	var openPortList []string
	var closedPortList []string

	if dnsScan(dstHostname, 53) {
		openPortList = append(openPortList, "53")
	} else {
		closedPortList = append(closedPortList, "53")
	}

	if snmpScan(dstHostname, 161) {
		openPortList = append(openPortList, "161")
	} else {
		closedPortList = append(closedPortList, "161")
	}

	printResults(&openPortList, &closedPortList)
}

// Since DNS is a UDP process, we need to send a valid DNS request to illicit a response from a service
// A response even if empty signifies presence of a DNS service
func dnsScan(dstHostname string, dstPort uint64) bool {
	target := "google.com"
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(target+".", dns.TypeA)

	log.Printf("[UDP Port %d] Sending DNS request to %s", dstPort, dstHostname)
	_, _, err := c.Exchange(&m, dstHostname+":"+strconv.FormatUint(dstPort, 10))
	if err != nil {
		log.Printf("[UDP Port %d] %v", dstPort, err)
		log.Printf("[UDP Port %d] is CLOSED", dstPort)
		return false
	} else {
		log.Printf("[UDP Port %d] DNS Response received", dstPort)
		log.Printf("[UDP Port %d] DNS is OPEN", dstPort)
		return true
	}

	// No need to do anything with the DNS Response
}

func snmpScan(dstHostname string, dstPort uint64) bool {
	gosnmp.Default.Target = dstHostname
	gosnmp.Default.Connect()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"} // we dont really care about what the query result is
	_, err := gosnmp.Default.Get(oids)                         // Get() accepts up to g.MAX_OIDS

	if err != nil {
		log.Printf("[UDP Port %d] %v", dstPort, err)
		log.Printf("[UDP Port %d] is CLOSED", dstPort)
		return false
	} else {
		log.Printf("[UDP Port %d] SNMP Response received", dstPort)
		log.Printf("[UDP Port %d] SNMP is OPEN", dstPort)
		return true
	}

}

func printResults(openPortList *[]string, closedPortList *[]string) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Port", "Protocol", "Status"})

	for _, op := range *openPortList {
		table.Append([]string{op, "UDP", "OPEN"})
	}
	for _, cp := range *closedPortList {
		table.Append([]string{cp, "UDP", "CLOSED"})
	}

	table.Render()
}

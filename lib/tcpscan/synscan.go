package tcpscan

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/bongtavas/gonmap/lib/netutil"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/olekukonko/tablewriter"
)

// PortStatus describes the status of the port
type PortStatus struct {
	port   uint64
	status string
}

// SynScan determines open and close ports in host using TCP SYN Scanning
func SynScan(dstHostname string, dstPortList []uint64) {
	portStatusChan := make(chan PortStatus)

	dstIP := netutil.ResolveHostname(dstHostname)

	var portStatusList []PortStatus

	for _, dstPort := range dstPortList {
		log.Printf("Scanning %s:%d", dstIP, dstPort)
		go scanPort(dstIP, dstPort, portStatusChan)
	}

	for i := 0; i < len(dstPortList); i++ {
		portStatus := <-portStatusChan
		portStatusList = append(portStatusList, portStatus)
	}

	printResults(&dstIP, &portStatusList)
}

func scanPort(dstIP net.IP, dstPortInt uint64, c chan PortStatus) {
	dstPort := layers.TCPPort(dstPortInt)
	srcIP, srcPortInt := netutil.GetLocalIPPort(dstIP)
	srcPort := layers.TCPPort(srcPortInt)
	tcpHeader := makeTCPHeader(srcPort, dstPort)
	ipHeader := makeIPHeader(srcIP, dstIP) // this for checksum

	tcpHeader.SetNetworkLayerForChecksum(ipHeader)

	// We only serialize the TCP layer because the
	// socket we get with net.ListenPacket already wraps our data in IPv4 packets.
	// We still need the IP layer to compute for checksums though
	buf := serializeTCPRequest(tcpHeader)

	// Sniff incoming packets
	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Write TCP SYN packet to the wire
	log.Printf("[Port %d] Writing TCP SYN Packet", dstPortInt)
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
		log.Fatalf("[Port %d] %v", dstPortInt, err)
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatalf("[Port %d] %v", dstPortInt, err)
	}

	for {
		b := make([]byte, 4096)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Printf("[Port %d] Error reading packet %s: ", dstPortInt, err)
			c <- PortStatus{dstPortInt, "CLOSED"}
			break
		} else if addr.String() == dstIP.String() {
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcPort {
					log.Printf("[Port %d] TCP response received ", dstPortInt)
					if tcp.SYN && tcp.ACK {
						log.Printf("[Port %d] Received SYN and ACK ", dstPortInt)
						log.Printf("[Port %d] is OPEN ", dstPortInt)
						c <- PortStatus{dstPortInt, "OPEN"}
					} else {
						log.Printf("[Port %d] Did not receive SYN and ACK ", dstPortInt)
						log.Printf("[Port %d] is CLOSED ", dstPortInt)
						c <- PortStatus{dstPortInt, "CLOSED"}
					}
					break
				}
			}
		}
	}
}

func makeTCPHeader(srcPort layers.TCPPort, dstPort layers.TCPPort) *layers.TCP {
	tcpHeader := layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Seq:     1105024978,
		SYN:     true,
		Window:  14600,
	}
	return &tcpHeader
}

func makeIPHeader(srcIP net.IP, dstIP net.IP) *layers.IPv4 {
	ipHeader := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Protocol: layers.IPProtocolTCP,
	}
	return &ipHeader
}

func serializeTCPRequest(tcpHeader *layers.TCP) gopacket.SerializeBuffer {
	buf := gopacket.NewSerializeBuffer()

	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	if err := gopacket.SerializeLayers(buf, opts, tcpHeader); err != nil {
		log.Fatal(err)
	}
	return buf
}

func printResults(dstIP *net.IP, portStatusList *[]PortStatus) {

	table := tablewriter.NewWriter(os.Stdout)
	table.ClearRows()
	table.SetHeader([]string{"Port", "Protocol", "Status"})

	fmt.Printf("Results for : %s\n", dstIP.String())
	for _, portStatus := range *portStatusList {
		portStr := strconv.FormatUint(portStatus.port, 10)
		table.Append([]string{portStr, "TCP", portStatus.status})
	}

	table.Render()
}

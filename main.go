package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Get the local IP and port based on our destination IP
// This is done by sending a dummy UDP packet to our destination IP
func getLocalIPPort(dstIP net.IP) (net.IP, int) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstIP.String()+":12345")
	if err != nil {
		log.Fatal(err)
	}

	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port
		}
	}
	log.Fatal("could not get local ip: " + err.Error())
	return nil, -1
}

func resolveHostname(hostname string) net.IP {
	log.Printf("Resolving %s", hostname)
	dstaddrs, err := net.LookupIP(hostname)
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
	resolvedIP := dstaddrs[0]
	log.Printf("Resolved IP: %s", resolvedIP)
	return resolvedIP.To4()
}

func wrapDstPort(dstPortStr string) layers.TCPPort {
	var dstport layers.TCPPort
	if d, err := strconv.ParseUint(os.Args[2], 10, 16); err != nil {
		log.Fatal(err)
	} else {
		dstport = layers.TCPPort(d)
	}
	return dstport
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

func main() {
	if len(os.Args) != 3 {
		log.Printf("Usage: %s <host/ip> <port>\n", os.Args[0])
		os.Exit(-1)
	}
	log.Printf("gonmap %s %s", os.Args[1], os.Args[2])

	dstHostname := os.Args[1]
	dstPortStr := os.Args[2]

	dstIP := resolveHostname(dstHostname)
	dstPort := wrapDstPort(dstPortStr)

	srcIP, srcPortInt := getLocalIPPort(dstIP)
	srcPort := layers.TCPPort(srcPortInt)

	log.Printf("Source IP: %v", srcIP.String())
	log.Printf("Source Port: %v", srcPort.String())
	log.Printf("Destination IP: %v", dstIP.String())
	log.Printf("Destination Port: %v", dstPort.String())

	tcpHeader := makeTCPHeader(srcPort, dstPort)
	ipHeader := makeIPHeader(srcIP, dstIP) // this for checksum

	tcpHeader.SetNetworkLayerForChecksum(ipHeader)

	// We only serialize the TCP layer because the
	// socket we get with net.ListenPacket already wraps our data in IPv4 packets.
	// We still need the IP layer to compute for checksums though
	buf := serializeTCPRequest(tcpHeader)

	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Writing TCP SYN Packet")
	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstIP}); err != nil {
		log.Fatal(err)
	}

	// Set deadline so we don't wait forever.
	if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, 4096)
		log.Println("Reading from socket")
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Println("Error reading packet: ", err)
			return
		} else if addr.String() == dstIP.String() {
			// Decode a packet
			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
			// Get the TCP layer from this packet
			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
				tcp, _ := tcpLayer.(*layers.TCP)

				if tcp.DstPort == srcPort {
					log.Printf("TCP response received")
					if tcp.SYN && tcp.ACK {
						log.Printf("Received SYN and ACK")
						log.Printf("Port %d is OPEN\n", dstPort)
					} else {
						log.Printf("Port %d is CLOSED\n", dstPort)
					}
					return
				}
			}
		} else {
			log.Printf("Got packet not matching address")
		}
	}
}

# gonmap
nmap written in Go

## Features
 * [x] TCP SYN Scan
 * [x] UDP Scan  - Port 53 DNS, Port 161 SNMP (Enable through 'udp' flag)
 * [ ] Service/Banner Grabbing
 * [ ] OS Detection

## Building
```bash
go build -o gonmap main.go
```

## Running
Note: `sudo` is needed to be able to write directly to the wire
```bash
sudo ./gonmap -p <comma-separted-ports> <hostname>
```

Ex:
```bash
sudo ./gonmap google.com -p 80,443,22
```

You can also specify port range using the '-' delimiter
```bash
sudo ./gonmap -p 1-22,443-600,5443 google.com
```

If you don't specify a port argument, the default 1000 ports of nmap will be used.

UDP Scanning is disbled by default, use the flag "udp" to enable it
```bash
sudo ./gonmap -udp 8.8.8.8
```

```bash
sudo ./gonmap -udp demo.snmplabs.com
```

# gonmap
nmap written in Go

This is a submission of __Team KnoxNaman__ for Machine Problem 2 of EE 298-K Information Security class
- Romelio Tavas Jr. [@bongtavas](https://github.com/bongtavas) 
- Rene Josiah Quinto [@rjmquinto](https://github.com/rjmquinto) 
- Josen Daniel De Leon [@jddl68](https://github.com/jddl68) 

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
sudo ./gonmap -sS -sU -p <comma-separted-ports> <hostname>
```

<hostname> can be the following forms:
- Hostname e.g. pleni.upd.edu.ph
- IP Address e.g. 202.92.128.181
- IP range e.g. 202.92.128.1-254
- IP CIDR block e.g. 202.92.128.0/24

## TCP SYN Scanning
Use flag `-sS` to use TCP SYN scanning
Ex:
```bash
sudo ./gonmap -sS -p 80,443,22 google.com 
```

You can also specify port range using the '-' delimiter
```bash
sudo ./gonmap -sS -p 1-22,443-600,5443 google.com
```

If you don't specify a port argument, the default 1000 ports of nmap will be used.

## UDP Scanning
Use flag `-sU` to use UDP scanning
```bash
sudo ./gonmap -sU 8.8.8.8
```

```bash
sudo ./gonmap -sU demo.snmplabs.com
```

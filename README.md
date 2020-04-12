# gonmap
nmap written in Go

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

# gonmap
nmap written in Go

## Building
```bash
go build -o gonmap main.go
```

## Running
```bash
./gonmap <hostname> <port>
```

Example: Check if google.com port 80 is open
```bash
./gonmap google.com 80
```

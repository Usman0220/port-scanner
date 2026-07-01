# Port Scanner

A fast, concurrent port scanner written in Go. Scans random public IP addresses on a specified port with service identification.

## Features

- High-concurrency scanning with configurable worker count
- Banner grabbing and service identification (FTP, SSH, HTTP, MySQL, Redis, etc.)
- Automatic probe selection based on target port
- Real-time progress display
- Output to file

## Usage

```bash
# Build
go build -o port-scanner main.go

# Scan port 5678 with 500 workers, 10000 IPs
./port-scanner -port 5678 -w 500 -n 10000

# Scan with output file
./port-scanner -port 22 -w 1000 -n 50000 -o results.txt

# Custom timeout (5s)
./port-scanner -port 80 -timeout 5s -n 10000
```

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | 5678 | Port to scan |
| `-w` | 500 | Number of goroutine workers |
| `-timeout` | 2s | Connection timeout |
| `-n` | 10000 | Number of random IPs to scan |
| `-o` | "" | Output file path |

## Supported Services

FTP, SSH, Telnet, SMTP, POP3, IMAP, HTTP (Nginx/Apache/IIS/Cloudflare), MySQL, PostgreSQL, Redis, MongoDB, n8n, Grafana, Kibana, Jenkins, Git

## License

MIT

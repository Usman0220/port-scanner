<div align="center">

```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║   ░██████╗░██████╗░██████╗░██╗░░░██╗██╗██████╗░███████╗    ║
║   ██╔════╝██╔═══██╗██╔══██╗██║░░░██║██║██╔══██╗██╔════╝    ║
║   ╚█████╗░██║██╗██║██║░░██║██║░░░██║██║██████╔╝█████╗░░    ║
║   ░╚═══██╗╚██████╔╝██║░░██║██║░░░██║██║██╔══██╗██╔══╝░░    ║
║   ██████╔╝░╚═██╔═╝░██████╔╝╚██████╔╝██║██║░░██║███████╗    ║
║   ╚═════╝░░░░╚═╝░░░╚═════╝░░╚═════╝░╚═╝╚═╝░░╚═╝╚══════╝  ║
║                                                              ║
║              [ CONCURRENT PORT SCANNER ]                      ║
║           Go × Nuclei × Service Identification               ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
```

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Nuclei](https://img.shields.io/badge/Nuclei-F74B03?style=for-the-badge&logo=projectdiscovery&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux-blue?style=for-the-badge)

**Lightning-fast concurrent port scanner with banner grabbing, service fingerprinting, and Nuclei template integration for automated vulnerability scanning.**

[Features](#features) • [Quick Start](#quick-start) • [Nuclei Integration](#nuclei-integration) • [Workflow Pipeline](#workflow-pipeline) • [Flags](#flags) • [Examples](#examples)

</div>

---

## Features

| Feature | Description |
|---------|-------------|
| 🚀 **High Concurrency** | Configurable goroutine pool — scan thousands of IPs simultaneously |
| 🔍 **Banner Grabbing** | Reads service banners for accurate fingerprinting |
| 🧠 **Smart Probing** | Auto-selects the right probe per port (HTTP, SMTP, Redis, etc.) |
| 🏷️ **Service Detection** | Identifies 15+ services: FTP, SSH, HTTP, MySQL, Redis, MongoDB, n8n, Grafana... |
| ⚡ **Real-time Progress** | Live stats: scanned / open / verified |
| 📦 **File Output** | Save results for pipeline consumption |
| 🎯 **Nuclei Ready** | Output format plugs directly into Nuclei for vuln scanning |

---

## Quick Start

```bash
# Clone
git clone https://github.com/Usman0220/port-scanner.git
cd port-scanner

# Build
go build -o port-scanner main.go

# Run — scan port 5678 with 500 workers, 10k random IPs
./port-scanner -port 5678 -w 500 -n 10000
```

---

## Nuclei Integration

The scanner outputs results in a format Nuclei can consume directly. Pipe open targets into Nuclei for automated vulnerability detection.

### Scan & Pipe to Nuclei

```bash
# Step 1: Find open FTP servers
./port-scanner -port 21 -w 1000 -n 50000 -o ftp-open.txt

# Step 2: Extract IPs only for Nuclei
awk -F'[|]' '{print $1}' ftp-open.txt | sed 's/\[OPEN\] //' | cut -d: -f1 > ftp-targets.txt

# Step 3: Run Nuclei with FTP templates
nuclei -l ftp-targets.txt -tags ftp -t ~/.local/nuclei-templates/
```

### Scan Multiple Ports → Nuclei

```bash
#!/bin/bash
# multi-port-scan.sh — scan ports, then audit each service with Nuclei

PORTS=(21 22 80 443 3306 5432 6379 8080 8443 9090 27017)
WORKERS=1000
IPS=20000

for port in "${PORTS[@]}"; do
    echo "[*] Scanning port $port..."
    ./port-scanner -port $port -w $WORKERS -n $IPS -o "results-port${port}.txt"
    
    # Extract live hosts
    awk -F'[|]' '{print $1}' "results-port${port}.txt" | \
        sed 's/\[OPEN\] //' | cut -d: -f1 | sort -u > "targets-port${port}.txt"
    
    count=$(wc -l < "targets-port${port}.txt")
    echo "[+] Found $count live hosts on port $port"
    
    if [ "$count" -gt 0 ]; then
        echo "[*] Running Nuclei against port $port targets..."
        nuclei -l "targets-port${port}.txt" -p-port $port \
            -t ~/.local/nuclei-templates/ -severity critical,high,medium \
            -o "nuclei-port${port}.txt" -silent
    fi
done

echo "[+] All done. Check nuclei-port*.txt for findings."
```

### Nuclei with Specific Templates

```bash
# FTP anonymous access, brute force, known CVEs
nuclei -l ftp-targets.txt -tags ftp -severity critical,high

# SSH vulnerabilities
nuclei -l ssh-targets.txt -tags ssh -t ~/.local/nuclei-templates/

# HTTP vuln scan — take full advantage of templates
nuclei -l http-targets.txt -tags http -severity critical,high,medium

# Redis unauthorized access
nuclei -l redis-targets.txt -tags redis

# MongoDB no-auth
nuclei -l mongo-targets.txt -tags mongodb

# Everything — full audit
nuclei -l targets.txt -t ~/.local/nuclei-templates/ -severity critical,high,medium,low
```

---

## Workflow Pipeline

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   PORT SCANNER  │────▶│  TARGET FILTER   │────▶│     NUCLEI      │
│                 │     │                  │     │                 │
│  Random IPs     │     │  awk / grep      │     │  Vuln Templates │
│  Banner Grab    │────▶│  Extract IPs     │────▶│  CVE Detection  │
│  Service ID     │     │  Deduplicate     │     │  Misconfigs     │
└─────────────────┘     └──────────────────┘     └─────────────────┘
        │                       │                        │
        ▼                       ▼                        ▼
   results.txt           targets.txt              findings.txt
```

---

## Examples

### Scan n8n instances on port 5678
```bash
./port-scanner -port 5678 -w 500 -n 30000 -o n8n-results.txt
nuclei -l n8n-targets.txt -tags n8n -severity critical,high
```

### Mass HTTP discovery
```bash
./port-scanner -port 80 -w 2000 -n 100000 -o http-results.txt
./port-scanner -port 443 -w 2000 -n 100000 -o https-results.txt
cat http-results.txt https-results.txt | awk -F'[|]' '{print $1}' | \
    sed 's/\[OPEN\] //' | cut -d: -f1 | sort -u > web-targets.txt
nuclei -l web-targets.txt -tags http -severity critical,high,medium -o web-findings.txt
```

### Find exposed databases
```bash
# MySQL
./port-scanner -port 3306 -w 500 -n 50000 -o mysql.txt
nuclei -l mysql-targets.txt -tags mysql

# Redis
./port-scanner -port 6379 -w 500 -n 50000 -o redis.txt
nuclei -l redis-targets.txt -tags redis

# MongoDB
./port-scanner -port 27017 -w 500 -n 50000 -o mongo.txt
nuclei -l mongo-targets.txt -tags mongodb

# PostgreSQL
./port-scanner -port 5432 -w 500 -n 50000 -o postgres.txt
nuclei -l postgres-targets.txt -tags postgresql
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `5678` | Port to scan |
| `-w` | `500` | Goroutine workers (higher = faster, more bandwidth) |
| `-timeout` | `2s` | Connection timeout per IP |
| `-n` | `10000` | Number of random public IPs to scan |
| `-o` | `""` | Output file path |

---

## Output Format

```
[OPEN] 192.168.1.100:21 | FTP | 220 ProFTPD 1.3.6 Server ready
[OPEN] 10.0.0.55:80 | HTTP/Apache | HTTP/1.1 200 OK Apache/2.4.52
[OPEN] 172.16.0.1:6379 | Redis | Redis server version 7.0.5
```

Each line: `[STATUS] IP:PORT | SERVICE | BANNER`

---

## Supported Services

<details>
<summary><b>Click to expand full list</b></summary>

| Port(s) | Service | Probe |
|---------|---------|-------|
| 21 | FTP | Banner |
| 22 | SSH | Banner |
| 23 | Telnet | Banner |
| 25, 587 | SMTP | `EHLO test` |
| 80, 8080, 8443, 443, 5678, 3000, 8000, 8888, 9090 | HTTP | `GET /` |
| 110 | POP3 | Banner |
| 143 | IMAP | Banner |
| 3306 | MySQL | Handshake |
| 5432 | PostgreSQL | Startup |
| 6379 | Redis | `INFO` |
| 27017 | MongoDB | Hello |
| Default | HTTP | `GET /` |

**Service fingerprinting:** FTP, SSH, Telnet, SMTP, POP3, IMAP, HTTP/Nginx, HTTP/Apache, HTTP/IIS, HTTP/Cloudflare, MySQL, PostgreSQL, Redis, MongoDB, n8n, Grafana, Kibana, Jenkins, Gitea/GitLab

</details>

---

## Requirements

- **Go** ≥ 1.20
- **Nuclei** (optional) — for vulnerability scanning: `go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest`
- **Nuclei templates** — `nuclei -update-templates`

---

## Disclaimer

This tool is for **authorized security testing and research purposes only**. Scanning networks without permission is illegal. Use responsibly.

---

<div align="center">

**Built with Go** | **Powered by Nuclei Templates**

```
  /\_/\    Port Scanner
 ( o.o )   + Nuclei
  > ^ <
```

</div>

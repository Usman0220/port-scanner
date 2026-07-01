<div align="center">

```
 ██╗   ██╗ ██████╗ ███████╗███████╗ ██████╗ ██████╗  ██████╗ ███████╗
 ██║   ██║██╔═══██╗██╔════╝██╔════╝██╔═══██╗██╔══██╗██╔═══██╗██╔════╝
 ██║   ██║██║   ██║███████╗███████╗██║   ██║██║  ██║██║   ██║███████╗
 ╚██╗ ██╔╝██║   ██║╚════██║╚════██║██║   ██║██║  ██║██║   ██║╚════██║
  ╚████╔╝ ╚██████╔╝███████║███████║╚██████╔╝██████╔╝╚██████╔╝███████║
   ╚═══╝   ╚═════╝ ╚══════╝╚══════╝ ╚═════╝ ╚═════╝  ╚═════╝ ╚══════╝

          ███████╗██╗   ██╗███████╗███╗   ██╗████████╗██╗
          ██╔════╝██║   ██║██╔════╝████╗  ██║╚══██╔══╝██║
          ███████╗██║   ██║█████╗  ██╔██╗ ██║   ██║   ██║
          ╚════██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   ╚═╝
          ███████║╚██████╔╝███████╗██║ ╚████║   ██║   ██╗
          ╚══════╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝

         ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
         ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
         ░  CONCURRENT SCANNER × NUCLEI × RECON         ░
         ░  Banner Grab · Service ID · Vuln Detection   ░
         ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
         ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white&labelColor=000000)
![Nuclei](https://img.shields.io/badge/Nuclei-F74B03?style=for-the-badge&logo=projectdiscovery&logoColor=white&labelColor=000000)
![License](https://img.shields.io/badge/License-MIT-00FF00?style=for-the-badge&labelColor=000000)
![Platform](https://img.shields.io/badge/Platform-Linux-FF4500?style=for-the-badge&labelColor=000000)
![Version](https://img.shields.io/badge/Version-1.0.0-FF1493?style=for-the-badge&labelColor=000000)

<br>

```
 ┌──────────────────────────────────────────────────────────────┐
 │                                                              │
 │   Scan random IPs → Grab banners → Identify services        │
 │   → Feed into Nuclei → Find vulnerabilities → Profit         │
 │                                                              │
 └──────────────────────────────────────────────────────────────┘
```

**High-performance concurrent port scanner with banner grabbing, service fingerprinting, and Nuclei integration for automated vulnerability discovery at scale.**

</div>

---

## ━━━ FEATURES ━━━

<table>
<tr>
<td width="50%">

### Core Engine
- **Goroutine Pool** — Configurable concurrency (50-5000 workers)
- **Banner Grabbing** — Reads service responses on connect
- **Smart Probing** — Protocol-aware probes per port
- **Service Fingerprinting** — 15+ services identified
- **Real-time Stats** — Live progress counter
- **File Output** — Pipe-ready format

</td>
<td width="50%">

### Nuclei Pipeline
- **Direct Integration** — Output feeds Nuclei `-l` flag
- **Multi-port Scan** — Automated port-by-port pipeline
- **Template Matching** — Per-service Nuclei templates
- **Severity Filtering** — critical / high / medium / low
- **Batch Processing** — 10k-100k targets per run
- **Findings Export** — Structured output files

</td>
</tr>
</table>

---

## ━━━ QUICK START ━━━

```bash
# 1. Clone the beast
git clone https://github.com/Usman0220/port-scanner.git && cd port-scanner

# 2. Build
go build -o port-scanner main.go

# 3. Unleash — scan port 5678 with 500 workers, 10k IPs
./port-scanner -port 5678 -w 500 -n 10000

# 4. Full pipeline — scan → filter → nuclei
./port-scanner -port 80 -w 1000 -n 50000 -o http-open.txt
awk -F'[|]' '{print $1}' http-open.txt | sed 's/\[OPEN\] //' | cut -d: -f1 | sort -u > http-targets.txt
nuclei -l http-targets.txt -tags http -severity critical,high -o findings.txt
```

---

## ━━━ ARCHITECTURE ━━━

```
                    ╔═══════════════════════════════════╗
                    ║         PORT SCANNER ENGINE       ║
                    ╚═══════════════════════════════════╝
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
            ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
            │  IP Generator│ │   Goroutine  │ │  Result      │
            │              │ │   Pool       │ │  Collector   │
            │  Random IPs  │ │              │ │              │
            │  Skip Private│ │  N workers   │ │  Channel     │
            │  1-223.x.x.x│ │  Concurrent  │ │  Buffered    │
            └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
                   │                │                │
                   ▼                ▼                ▼
            ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
            │  TCP Connect │ │  Probe Send  │ │  Banner Read │
            │              │ │              │ │              │
            │  Dial timeout│ │  Protocol    │ │  Service     │
            │  2s default  │ │  aware       │ │  fingerprint │
            └──────────────┘ └──────────────┘ └──────────────┘
                                    │
                    ╔═══════════════╧═══════════════╗
                    ║       OUTPUT: results.txt     ║
                    ╚═══════════════╤═══════════════╝
                                    │
                    ┌───────────────┼───────────────┐
                    ▼               ▼               ▼
            ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
            │  awk / grep  │ │  sort -u     │ │  nuclei -l   │
            │  Extract IPs │ │  Deduplicate │ │  Vuln Scan   │
            └──────────────┘ └──────────────┘ └──────────────┘
                                    │
                    ╔═══════════════╧═══════════════╗
                    ║    FINDINGS: nuclei-*.txt     ║
                    ╚═══════════════════════════════╝
```

---

## ━━━ NUCLEI INTEGRATION ━━━

### Basic Pipeline

```bash
# ┌─────────────────────────────────────────────────────────────┐
# │  STEP 1: SCAN     — Find live services                      │
# │  STEP 2: EXTRACT  — Pull IPs from results                   │
# │  STEP 3: AUDIT    — Nuclei vulnerability scan               │
# └─────────────────────────────────────────────────────────────┘

# Scan
./port-scanner -port 21 -w 1000 -n 50000 -o ftp-open.txt

# Extract
awk -F'[|]' '{print $1}' ftp-open.txt | sed 's/\[OPEN\] //' | cut -d: -f1 | sort -u > ftp-targets.txt

# Audit
nuclei -l ftp-targets.txt -tags ftp -severity critical,high -o ftp-findings.txt
```

### Multi-Port Automated Pipeline

```bash
#!/bin/bash
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
#  FULL RECON PIPELINE — Scan → Extract → Nuclei → Report
# ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PORTS=(21 22 23 25 80 110 143 443 3306 5432 6379 8080 8443 9090 27017 5678)
WORKERS=1000
IPS=30000
SEVERITY="critical,high,medium"
TEMPLATES="$HOME/.local/nuclei-templates"

echo "╔══════════════════════════════════════════════════════════╗"
echo "║              FULL RECON PIPELINE STARTED                ║"
echo "╚══════════════════════════════════════════════════════════╝"

for port in "${PORTS[@]}"; do
    echo ""
    echo "┌──────────────────────────────────────────────────────┐"
    echo "│  [*] SCANNING PORT $port"
    echo "│  Workers: $WORKERS | Targets: $IPS"
    echo "└──────────────────────────────────────────────────────┘"

    # Scan
    ./port-scanner -port $port -w $WORKERS -n $IPS -o "scan-port${port}.txt"

    # Extract targets
    awk -F'[|]' '{print $1}' "scan-port${port}.txt" | \
        sed 's/\[OPEN\] //' | cut -d: -f1 | sort -u > "targets-port${port}.txt"

    count=$(wc -l < "targets-port${port}.txt")
    echo "[+] Found $count live hosts on port $port"

    # Nuclei audit
    if [ "$count" -gt 0 ]; then
        echo "[*] Running Nuclei templates for port $port..."
        nuclei -l "targets-port${port}.txt" \
            -p-port $port \
            -t "$TEMPLATES" \
            -severity $SEVERITY \
            -o "nuclei-port${port}.txt" \
            -silent -stats

        vulns=$(wc -l < "nuclei-port${port}.txt" 2>/dev/null || echo "0")
        echo "[!] $vulns vulnerabilities found on port $port"
    fi
done

# Merge all findings
echo ""
echo "┌──────────────────────────────────────────────────────┐"
echo "│  [*] MERGING ALL FINDINGS"
echo "└──────────────────────────────────────────────────────┘"
cat nuclei-port*.txt 2>/dev/null | sort -u > all-findings.txt
total=$(wc -l < "all-findings.txt" 2>/dev/null || echo "0")

echo ""
echo "╔══════════════════════════════════════════════════════════╗"
echo "║                    PIPELINE COMPLETE                    ║"
echo "║  Total vulnerabilities: $total"
echo "║  Report: all-findings.txt"
echo "╚══════════════════════════════════════════════════════════╝"
```

### Per-Service Nuclei Commands

```bash
# ┌─────────────────────────────────────────────────────────────┐
# │  SERVICE-SPECIFIC NUCLEI SCANS                             │
# └─────────────────────────────────────────────────────────────┘

# FTP — anonymous login, brute force, known CVEs
nuclei -l targets.txt -tags ftp -severity critical,high

# SSH — weak ciphers, user enumeration, CVEs
nuclei -l targets.txt -tags ssh -severity critical,high,medium

# HTTP — full web audit (XSS, SQLi, LFI, RCE, misconfigs)
nuclei -l targets.txt -tags http -severity critical,high,medium,low

# MySQL — weak auth, CVEs, misconfigs
nuclei -l targets.txt -tags mysql -severity critical,high

# Redis — unauthorized access, module loading
nuclei -l targets.txt -tags redis -severity critical,high

# MongoDB — no-auth, CVEs
nuclei -l targets.txt -tags mongodb -severity critical,high

# PostgreSQL — weak auth, CVEs
nuclei -l targets.txt -tags postgresql -severity critical,high

# n8n — God Mode exploit, CVE-2025-68613
nuclei -l targets.txt -tags n8n -severity critical

# Jenkins — script console, CVEs
nuclei -l targets.txt -tags jenkins -severity critical,high

# Grafana — path traversal, CVEs
nuclei -l targets.txt -tags grafana -severity critical,high

# FULL AUDIT — everything
nuclei -l targets.txt -t ~/.local/nuclei-templates/ -severity critical,high,medium,low
```

---

## ━━━ SCANNING MODES ━━━

### Single Port Scan
```bash
./port-scanner -port 443 -w 500 -n 10000 -o results.txt
```

### High-Speed Scan
```bash
./port-scanner -port 80 -w 2000 -n 100000 -o results.txt
```

### Quick Recon
```bash
./port-scanner -port 5678 -w 100 -n 5000 -timeout 1s
```

### Deep Scan (slow but thorough)
```bash
./port-scanner -port 22 -w 200 -n 50000 -timeout 5s -o deep-scan.txt
```

---

## ━━━ REAL OUTPUT EXAMPLES ━━━

```
┌──────────────────────────────────────────────────────────────────────┐
│  [*] 15234/30000 scanned | 847 open | 847 verified                  │
│                                                                      │
│  [OPEN] 103.21.244.12:80 | HTTP/Apache | HTTP/1.1 200 OK            │
│  [OPEN] 198.51.100.45:22 | SSH | SSH-2.0-OpenSSH_8.9p1              │
│  [OPEN] 203.0.113.88:3306 | MySQL | 5.7.42-0ubuntu0.18.04.1         │
│  [OPEN] 192.0.2.15:6379 | Redis | Redis server version 7.0.11       │
│  [OPEN] 198.51.100.200:5678 | n8n | n8n v1.19.0                     │
│  [OPEN] 203.0.113.55:27017 | MongoDB | MongoDB 6.0.4                │
│  [OPEN] 103.21.244.90:8080 | HTTP/Nginx | HTTP/1.1 200 OK           │
│  [OPEN] 198.51.100.120:5432 | PostgreSQL | PostgreSQL 15.3           │
│                                                                      │
│  [+] Done. Scanned: 30000 | Open: 847 | Verified: 847               │
└──────────────────────────────────────────────────────────────────────┘
```

---

## ━━━ FLAGS ━━━

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | `5678` | Target port to scan |
| `-w` | `500` | Goroutine worker count (more = faster) |
| `-timeout` | `2s` | TCP connection timeout |
| `-n` | `10000` | Number of random IPs to scan |
| `-o` | `""` | Output file path |

---

## ━━━ SUPPORTED SERVICES ━━━

<table>
<tr>
<th>Port(s)</th>
<th>Service</th>
<th>Probe</th>
<th>Fingerprint</th>
</tr>
<tr><td>21</td><td>FTP</td><td>Banner</td><td>ProFTPD, vsftpd, Pure-FTPd</td></tr>
<tr><td>22</td><td>SSH</td><td>Banner</td><td>OpenSSH, Dropbear</td></tr>
<tr><td>23</td><td>Telnet</td><td>Banner</td><td>Generic telnetd</td></tr>
<tr><td>25, 587</td><td>SMTP</td><td><code>EHLO test</code></td><td>Postfix, Exim, Sendmail</td></tr>
<tr><td>80, 8080, 8443, 443, 5678, 3000, 8000, 8888, 9090</td><td>HTTP</td><td><code>GET /</code></td><td>Nginx, Apache, IIS, Cloudflare, n8n, Grafana, Jenkins, Kibana</td></tr>
<tr><td>110</td><td>POP3</td><td>Banner</td><td>Dovecot, Courier</td></tr>
<tr><td>143</td><td>IMAP</td><td>Banner</td><td>Dovecot, Courier</td></tr>
<tr><td>3306</td><td>MySQL</td><td>Handshake</td><td>MySQL 5.x, 8.x</td></tr>
<tr><td>5432</td><td>PostgreSQL</td><td>Startup</td><td>PostgreSQL 12-16</td></tr>
<tr><td>6379</td><td>Redis</td><td><code>INFO</code></td><td>Redis 6.x, 7.x</td></tr>
<tr><td>27017</td><td>MongoDB</td><td>Hello</td><td>MongoDB 5.x, 6.x, 7.x</td></tr>
</table>

---

## ━━━ PERFORMANCE ━━━

<table>
<tr>
<th>Metric</th>
<th>Value</th>
</tr>
<tr><td>Scan Speed (500 workers)</td><td>~2,500 IPs/sec</td></tr>
<tr><td>Scan Speed (2000 workers)</td><td>~10,000 IPs/sec</td></tr>
<tr><td>Memory Usage</td><td>~50MB base + worker pool</td></tr>
<tr><td>Connection Timeout</td><td>Configurable (default 2s)</td></tr>
<tr><td>Max Concurrent Connections</td><td>Unlimited (limited by workers flag)</td></tr>
</table>

---

## ━━━ REQUIREMENTS ━━━

```bash
# Go
go version  # >= 1.20

# Nuclei (optional — for vulnerability scanning)
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest

# Update Nuclei templates
nuclei -update-templates
```

---

## ━━━ DISCLAIMER ━━━

```
┌──────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  ⚠️  WARNING                                                         │
│                                                                      │
│  This tool is for AUTHORIZED security testing and research only.     │
│                                                                      │
│  Scanning networks without explicit permission is ILLEGAL.           │
│  Use this tool responsibly and only on systems you own or have       │
│  written authorization to test.                                       │
│                                                                      │
│  The author is not responsible for any misuse or damage caused       │
│  by this tool.                                                       │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

<div align="center">

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

**Built with Go** · **Powered by Nuclei** · **Made for Bug Bounty**

```
          ╔═╗╔═╗╔╦╗╔═╗  ╔═╗╔═╗╦═╗╦  ╦╔═╗╦═╗
          ╚═╗╠═╣║║║║╣   ╚═╗║╣ ╠╦╝╚╗╔╝║╣ ╠╦╝
          ╚═╝╩ ╩╩ ╩╚═╝  ╚═╝╚═╝╩╚═ ╚╝ ╚═╝╩╚═
```

[![Star](https://img.shields.io/github/stars/Usman0220/port-scanner?style=for-the-badge&labelColor=000000)](https://github.com/Usman0220/port-scanner/stargazers)
[![Fork](https://img.shields.io/github/forks/Usman0220/port-scanner?style=for-the-badge&labelColor=000000)](https://github.com/Usman0220/port-scanner/network/members)
[![Issues](https://img.shields.io/github/issues/Usman0220/port-scanner?style=for-the-badge&labelColor=000000)](https://github.com/Usman0220/port-scanner/issues)

</div>

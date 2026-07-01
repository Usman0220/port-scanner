package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func isPrivateFast(b0, b1 byte) bool {
	if b0 == 0 || b0 == 127 || b0 == 224 || b0 == 240 {
		return true
	}
	if b0 == 10 || b0 == 192 && b1 == 168 {
		return true
	}
	if b0 == 100 && b1 >= 64 && b1 <= 127 {
		return true
	}
	if b0 == 169 && b1 == 254 {
		return true
	}
	if b0 == 172 && b1 >= 16 && b1 <= 31 {
		return true
	}
	if b0 == 198 && b1 == 18 {
		return true
	}
	return false
}

func randomIP() (byte, byte, byte, byte) {
	for {
		b0 := byte(rand.Intn(223) + 1)
		b1 := byte(rand.Intn(256))
		if !isPrivateFast(b0, b1) {
			return b0, b1, byte(rand.Intn(256)), byte(rand.Intn(254) + 1)
		}
	}
}

// getProbe returns the right probe for a given port
func getProbe(addr string, port int) []byte {
	switch port {
	case 21:
		// FTP just needs to wait for banner
		return nil
	case 22:
		return nil // SSH sends banner on connect
	case 23:
		return nil // Telnet sends banner
	case 25, 587:
		return []byte("EHLO test\r\n")
	case 80, 8080, 8443, 443, 5678, 3000, 8000, 8888, 9090:
		return []byte("GET / HTTP/1.0\r\nHost: " + addr + "\r\n\r\n")
	case 110:
		return nil // POP3 sends banner
	case 143:
		return nil // IMAP sends banner
	case 3306:
		return nil // MySQL sends handshake
	case 5432:
		return nil // PostgreSQL sends startup
	case 6379:
		return []byte("INFO\r\n")
	case 27017:
		return nil // MongoDB sends hello
	default:
		return []byte("GET / HTTP/1.0\r\n\r\n")
	}
}

// identifyService detects service from banner/response
func identifyService(port int, resp string) string {
	lower := strings.ToLower(resp)
	switch {
	case port == 21 || strings.Contains(lower, "220") && strings.Contains(lower, "ftp"):
		return "FTP"
	case port == 22 || strings.Contains(lower, "ssh-"):
		return "SSH"
	case port == 23 || strings.Contains(lower, "login:") || strings.Contains(lower, "password:"):
		return "Telnet"
	case strings.Contains(lower, "220 ") && (strings.Contains(lower, "smtp") || strings.Contains(lower, "mail")):
		return "SMTP"
	case strings.Contains(lower, "+ok"):
		return "POP3"
	case strings.Contains(lower, "* ok") || strings.Contains(lower, "* preauth"):
		return "IMAP"
	case strings.Contains(lower, "http/"):
		return "HTTP"
	case strings.Contains(lower, "mysql"):
		return "MySQL"
	case strings.Contains(lower, "redis"):
		return "Redis"
	case strings.Contains(lower, "mongo"):
		return "MongoDB"
	case strings.Contains(lower, "postgresql") || strings.Contains(lower, "pg_"):
		return "PostgreSQL"
	case strings.Contains(lower, "openresty") || strings.Contains(lower, "nginx"):
		return "HTTP/Nginx"
	case strings.Contains(lower, "apache"):
		return "HTTP/Apache"
	case strings.Contains(lower, "cloudflare"):
		return "HTTP/Cloudflare"
	case strings.Contains(lower, "iis"):
		return "HTTP/IIS"
	case strings.Contains(lower, "n8n"):
		return "n8n"
	case strings.Contains(lower, "grafana"):
		return "Grafana"
	case strings.Contains(lower, "kibana"):
		return "Kibana"
	case strings.Contains(lower, "jenkins"):
		return "Jenkins"
	case strings.Contains(lower, "gitea") || strings.Contains(lower, "gitlab"):
		return "Git"
	case strings.Contains(lower, "503 service unavailable"):
		return "HTTP/503"
	default:
		return "Unknown"
	}
}

func main() {
	port := flag.Int("port", 5678, "Port to scan")
	workers := flag.Int("w", 500, "Goroutine workers")
	timeout := flag.Duration("timeout", 2*time.Second, "Connection timeout")
	total := flag.Int("n", 10000, "Number of IPs to scan")
	output := flag.String("o", "", "Output file")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var opened, scanned, verified atomic.Int64
	var wg sync.WaitGroup
	results := make(chan string, 1024)

	var outf *os.File
	if *output != "" {
		var err error
		outf, err = os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer outf.Close()
	}

	go func() {
		for r := range results {
			fmt.Println(r)
			if outf != nil {
				outf.WriteString(r + "\n")
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				s := scanned.Load()
				v := verified.Load()
				o := opened.Load()
				fmt.Fprintf(os.Stderr, "\r[*] %d/%d scanned | %d open | %d verified", s, *total, o, v)
			case <-done:
				return
			}
		}
	}()

	ipChan := make(chan [4]byte, *workers*4)
	go func() {
		for i := 0; i < *total; i++ {
			b0, b1, b2, b3 := randomIP()
			ipChan <- [4]byte{b0, b1, b2, b3}
		}
		close(ipChan)
	}()

	dialer := net.Dialer{Timeout: *timeout}

	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ip := range ipChan {
				addr := fmt.Sprintf("%d.%d.%d.%d:%d", ip[0], ip[1], ip[2], ip[3], *port)
				conn, err := dialer.Dial("tcp", addr)
				scanned.Add(1)
				if err != nil {
					continue
				}
				opened.Add(1)

				conn.SetReadDeadline(time.Now().Add(3 * time.Second))

				// send probe if needed
				probe := getProbe(addr, *port)
				if probe != nil {
					conn.Write(probe)
				}

				// read response
				buf := make([]byte, 4096)
				n, err := conn.Read(buf)
				conn.Close()

				var resp string
				if err == nil && n > 0 {
					resp = string(buf[:n])
					if len(resp) > 200 {
						resp = resp[:200]
					}
					resp = strings.ReplaceAll(resp, "\r\n", " ")
					resp = strings.ReplaceAll(resp, "\n", " ")
				}

				if resp == "" {
					continue
				}

				service := identifyService(*port, resp)
				verified.Add(1)
				results <- fmt.Sprintf("[OPEN] %d.%d.%d.%d:%d | %s | %s", ip[0], ip[1], ip[2], ip[3], *port, service, resp)
			}
		}()
	}

	wg.Wait()
	close(done)
	close(results)

	fmt.Fprintf(os.Stderr, "\n[+] Done. Scanned: %d | Open: %d | Verified: %d\n", scanned.Load(), opened.Load(), verified.Load())
}

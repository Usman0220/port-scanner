package main

import (
	"context"
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

func main() {
	workers := flag.Int("w", 500, "Goroutine workers")
	timeout := flag.Duration("timeout", 2*time.Second, "Connection timeout")
	total := flag.Int("n", 10000, "Number of IPs to scan")
	stop := flag.Int("stop", 0, "Stop after finding this many FTP servers (0 = no limit)")
	output := flag.String("o", "ftp-targets.txt", "Output file for nuclei -l")
	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())

	var scanned, found atomic.Int64
	var wg sync.WaitGroup
	var mu sync.Mutex

	outf, err := os.Create(*output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer outf.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				s := scanned.Load()
				f := found.Load()
				fmt.Fprintf(os.Stderr, "\r[*] %d/%d scanned | %d FTP found", s, *total, f)
			case <-ctx.Done():
				return
			}
		}
	}()

	ipChan := make(chan [4]byte, *workers*4)
	go func() {
		for i := 0; i < *total; i++ {
			select {
			case <-ctx.Done():
				break
			default:
			}
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
				select {
				case <-ctx.Done():
					return
				default:
				}

				addr := fmt.Sprintf("%d.%d.%d.%d:21", ip[0], ip[1], ip[2], ip[3])
				conn, err := dialer.Dial("tcp", addr)
				scanned.Add(1)
				if err != nil {
					continue
				}

				// read FTP banner
				conn.SetReadDeadline(time.Now().Add(3 * time.Second))
				buf := make([]byte, 1024)
				n, err := conn.Read(buf)
				conn.Close()

				if err != nil || n == 0 {
					continue
				}

				banner := string(buf[:n])
				// verify it's FTP (220 response)
				if !strings.Contains(banner, "220") {
					continue
				}

			f := found.Add(1)
			ip := fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])

			mu.Lock()
			outf.WriteString(ip + "\n")
			outf.Sync()
			mu.Unlock()

			fmt.Println(ip)

				if *stop > 0 && int(f) >= *stop {
					cancel()
					return
				}
			}
		}()
	}

	wg.Wait()
	cancel()

	fmt.Fprintf(os.Stderr, "\n[+] Done. Scanned: %d | FTP found: %d | Output: %s\n", scanned.Load(), found.Load(), *output)
}

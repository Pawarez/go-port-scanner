package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

type PortResult struct {
	Port   int    `json:"port"`
	State  string `json:"state"`
	Banner string `json:"banner"`
}

func scanPort(hostname string, ports <-chan int, wg *sync.WaitGroup, results *[]PortResult, mu *sync.Mutex) {
	defer wg.Done()

	for p := range ports {
		address := fmt.Sprintf("%s:%d", hostname, p)
		conn, err := net.DialTimeout("tcp", address, 1*time.Second)

		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 4096)
		n, _ := conn.Read(buffer)
		banner := string(buffer[:n])

		conn.Close()

		mu.Lock() // Lock to append openPorts to resolve data race

		*results = append(*results, PortResult{
			Port:   p,
			State:  "Open",
			Banner: banner,
		})

		mu.Unlock() // Unlock

	}

}

func main() {
	hostname := flag.String("host", "scanme.nmap.org", "Hostname or IP address to scan")

	startPort := flag.Int("start", 1, "Start Port")
	endPort := flag.Int("end", 1024, "End Port")
	jsonFilename := flag.String("json","","Output results to JSON file (e.g., result.json)")

	flag.Parse()

	fmt.Printf("Scanning %s from port %d to %d...\n", *hostname, *startPort, *endPort)
	start := time.Now()

	ports := make(chan int, 100)

	var wg sync.WaitGroup
	var results []PortResult
	var mu sync.Mutex

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go scanPort(*hostname, ports, &wg, &results, &mu)
	}

	for i := *startPort; i <= *endPort; i++ {
		ports <- i
	}

	close(ports)

	wg.Wait()

	elapsed := time.Since(start)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	fmt.Printf("\n--- Scan Results ---\n")
	fmt.Printf("Target: %s\n", *hostname)
	fmt.Printf("Time taken: %v\n", elapsed)

	if len(results) == 0 {
		fmt.Println("No open ports found.")
	} else {
		fmt.Printf("%-10s %-10s %s\n", "PORT", "STATE", "SERVICE/VERSION")
		fmt.Println("--------------------------------------------")
		for _, r := range results {
			if r.Banner == "" {
				r.Banner = "Unknown"
			}
			cleanBanner := ""
			for _, char := range r.Banner {
				if char >= 32 && char <= 126 { 
					cleanBanner += string(char)
				}
			}

			fmt.Printf("%-10d %-10s %s\n", r.Port, r.State, cleanBanner)
		}
	}

	if *jsonFilename != "" {
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			return
		}
		err = os.WriteFile(*jsonFilename, data, 0644)
		if err != nil {
			fmt.Println("Error writing file:", err)
			return
		}
		fmt.Printf("\n[+] Results saved to %s successfully!\n", *jsonFilename)
	}

}

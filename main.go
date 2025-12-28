package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Pawarez/goscanner/pkg/scanner"
	"os"
	"strings"
	"time"
)

func main() {
	targetInput := flag.String("target", "127.0.0.1", "Target IP or CIDR (e.g., '192.168.1.1' or '192.168.1.0/24')")
	portsInput := flag.String("ports", "80,443", "Ports to scan (e.g., '80,443,8000-8080')")
	jsonFilename := flag.String("json", "", "Output results to JSON file")
	flag.Parse()

	targetPorts, err := scanner.ParsePorts(*portsInput)
	if err != nil {
		fmt.Printf("Error parsing ports: %v\n", err)
		os.Exit(1)
	}
	
	var targets []string
	
	if strings.Contains(*targetInput,"/") {
		targets, err = scanner.ParseCIDR(*targetInput) 

		if err != nil {
			fmt.Printf("Error Parsing CIDR: ", err) 
			os.Exit(1)
		}
	} else {
		// 1 ip
		targets = []string{*targetInput}
	}
	
	
	fmt.Printf("Scanning %d targets for %d ports...\n", len(targets), len(targetPorts))
	start := time.Now()

	var allResults []scanner.PortResult

	for _ , host := range targets {
		results := scanner.Run(host, targetPorts)

		allResults = append(allResults, results...)
	}

	elapsed := time.Since(start)

	fmt.Printf("\n--- Scan Results ---\n")
	fmt.Printf("Time taken: %v\n", elapsed)
	fmt.Printf("Total Open Ports Found: %d\n", len(allResults))

	if len(allResults) == 0 {
		fmt.Println("No open ports found.")
	} else {
		fmt.Println("----------------------------------------------------------------")
		fmt.Printf("%-15s %-10s %-10s %s\n", "HOST", "PORT", "STATE", "SERVICE/VERSION")
		fmt.Println("----------------------------------------------------------------")
		
		for _, r := range allResults {
			if r.Banner == "" {
				r.Banner = "Unknown"
			}
			cleanBanner := ""
			for _, char := range r.Banner {
				if char >= 32 && char <= 126 {
					cleanBanner += string(char)
				}
			}
			fmt.Printf("%-15s %-10d %-10s %s\n", r.Host, r.Port, r.State, cleanBanner)
		}
	}

	if *jsonFilename != "" {
		data, err := json.MarshalIndent(allResults, "", "  ")
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
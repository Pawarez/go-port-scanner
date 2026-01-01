#  High-Performance TCP Port Scanner (Go)

A fast, concurrent network reconnaissance tool written in **Go (Golang)**. Designed for speed and reliability, utilizing low-level concurrency patterns to scan network ranges efficiently without exhausting system resources.

## ⚠️ Legal Disclaimer
This tool is designed for **educational purposes and authorized security testing only**. 
The developer is not responsible for any misuse or damage caused by this tool. 
Please ensure you have permission before scanning any target network.

##  Key Features

- **High Performance:** Utilizes **Goroutines** to perform massive parallel scanning.
- **Resource Control:** Implements the **Semaphore pattern** to strictly limit concurrent threads (e.g., 100-1000 threads), preventing file descriptor exhaustion and network congestion.
- **Thread Safety:** Uses `sync.Mutex` to safely aggregate results during multi-threaded execution, preventing race conditions.
- **Smart Parsing:** Supports **CIDR notation** parsing (e.g., `192.168.1.0/24`) to scan entire subnets automatically.
- **Data Export:** Exports results to structured **JSON** format for integration with other security tools (e.g., SIEM, Log Analysis).

## Installation

Ensure you have [Go](https://go.dev/dl/) installed.

```bash
# Clone the repository
git clone https://github.com/yourusername/go-port-scanner.git

# Navigate to directory
cd go-port-scanner

# Run directly
go run main.go -h

# Example Command
go run main.go -target 192.168.1.1 -ports 1-1000
```

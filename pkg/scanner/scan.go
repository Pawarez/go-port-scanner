package scanner

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

type PortResult struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	State  string `json:"state"`
	Banner string `json:"banner"`
}


func Run(hostname string,portList []int) []PortResult {
	var results []PortResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	ports := make(chan int, 100) 

	for i := 0 ; i < 100; i++ {
		wg.Add(1)
		go worker(hostname,ports,&wg,&results,&mu) 

	}

	for _, p := range portList {
		ports <- p  
	}

	close(ports)

	wg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

func worker (hostname string,ports <- chan int,wg*sync.WaitGroup,results *[]PortResult,mu *sync.Mutex) {

	defer wg.Done()

	for p := range ports {
		address := fmt.Sprintf("%s:%d",hostname,p)

		conn,err := net.DialTimeout("tcp",address,1*time.Second)

		if err != nil {
			continue
		}

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 4096)
		n, _ := conn.Read(buffer)
		banner := string(buffer[:n])

		conn.Close()

		mu.Lock()
		*results = append(*results, PortResult{
			Host:   hostname,
			Port:   p,
			State:  "Open",
			Banner: banner,
		})
		mu.Unlock()
	}
}

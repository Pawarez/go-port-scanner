package scanner

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParsePorts(portRange string)  ([]int, error) {

	var ports []int

	tokens := strings.Split(portRange, ",") 

	for _, token := range tokens {
		token = strings.TrimSpace(token)

		if token == "" {
			continue
		}

		if strings.Contains(token,"-") {

			parts := strings.Split(token, "-")

			if len(parts) != 2 {
				return nil, fmt.Errorf("Wrong Format. %s",token)
			}

			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			
			if err1 != nil || err2 != nil {
				return nil, fmt.Errorf("invalid number in range: %s", token)
			}

			if start > end {
				return nil, fmt.Errorf("invalid range order: %d > %d", start, end)
			}

			for p := start; p <= end; p++ {
				ports = append(ports, p)
			}
		} else {

			port, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", token)
			}
			ports = append(ports, port)
		}
	} 


	return ports,nil 
}

func inc (ip net.IP) {
	for j := len(ip) - 1; j >= 0 ; j-- {
		ip[j]++
		if ip[j] > 0{
			break
		}
	}
}

func ParseCIDR(cidr string) ([]string, error) {

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}

	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil 
	}

	return ips, nil
}
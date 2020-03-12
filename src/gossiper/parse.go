package main

import "os"
import "net"
import "strings"
import "io/ioutil"

// Functions for parsing a target string as a path, CIDR range or single host.

func parse_target(target_str string) []string {
	//is the target a filepath that we can open?
	fd,err := os.Open(target_str)
	if err == nil {
		defer fd.Close()
		buf,err := ioutil.ReadAll(fd)
		if err == nil {
			output := []string{}
			parts := splitlines(string(buf))
			for _,part := range parts {
				if part != "" {
					output = append(output, part)
				}
			}
			return output
		}
	}
	//is the target a CIDR range we can expand?
	ip,ipnet,err := net.ParseCIDR(target_str)
	if err == nil {
		output := []string{}
		for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); increment_ip(ip) {
			output = append(output, ip.String())
		}
		return output
		
	}
	//treat the target as a single IP or hostname
	return append([]string{}, target_str)
}

func increment_ip(ip net.IP) {
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func splitlines(in_str string) []string {
	return strings.Split(strings.Replace(in_str, "\r\n", "\n", -1), "\n")
}



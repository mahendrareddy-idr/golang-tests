package main

import (
	"fmt"
	"regexp"
)

var ipv4Pattern = regexp.MustCompile(`^(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}$`)

func isValidIP(input string) bool {
	return ipv4Pattern.MatchString(input)
}

func main() {
	testIPs := []string{"192.168.1.1", "10.0.0.256", "256.0.0.1", "abc", "1.2.3.4"}

	for _, ip := range testIPs {
		fmt.Printf("%s -> %v\n", ip, isValidIP(ip))
	}
}

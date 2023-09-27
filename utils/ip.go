package utils

import "net"

func IPInCIDR(ipStr, cidrStr string) bool {
	// parse the IP address and CIDR block
	ip := net.ParseIP(ipStr)
	_, cidr, _ := net.ParseCIDR(cidrStr)

	// check if the IP address is within the CIDR block
	return cidr.Contains(ip)
}

type CIDRChecker struct {
	cidr *net.IPNet
}

func NewCIDRChecker(cidrStr string) *CIDRChecker {
	c := &CIDRChecker{}
	_, cidr, _ := net.ParseCIDR(cidrStr)
	c.cidr = cidr
	return c
}

func (self *CIDRChecker) Check(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return self.cidr.Contains(ip)
}

func (self *CIDRChecker) String() string {
	return self.cidr.String()
}

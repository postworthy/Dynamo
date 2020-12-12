package main

import (
	"encoding/json"
	"net"
	"strings"
)

type dnsResult struct {
	Domain string `json:"domain"`
	IP []net.IP   `json:"ips"`
	Error error   `json:error`
}

func (result * dnsResult) String() string{
	if result.Error != nil {
		return result.Error.Error()
	} else {
		var ips []string
		for _, ip := range result.IP {
			ips = append(ips, ip.String())
		}
		return strings.Join(ips, ",") + "," + result.Domain
	}
}

func (result * dnsResult) Json() string{
	json, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(json)
}
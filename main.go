package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var wgInner sync.WaitGroup

func main() {

	domainData, err := ioutil.ReadFile("domain-list.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	subdomainData, err := ioutil.ReadFile("subdomain-list.txt")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}

	domainList := SplitLines(string(domainData))
	subdomainList := SplitLines(string(subdomainData))
	results := make(chan dnsResult, 100)

	wg.Add(1)
	go printResults(results)

	for _, dns := range domainList {
		dns = strings.Trim(strings.Split(dns, "#")[0], " ")
		goroutines := maxParallelism()
		slices := int(math.Ceil(float64(len(subdomainList)) /  float64(goroutines)))
		for i := 0; i < slices; i++ {
			for _, username := range subdomainList[i:i+goroutines] {
				domain := username + "." + dns
				wg.Add(1)
				wgInner.Add(1)
				go lookupDomain(domain, results)
			}
			wgInner.Wait()
		}

	}

	wg.Wait()
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func printResults(results chan dnsResult) {
	defer wg.Done()
	for result := range results {
		if result.Error != nil {
			fmt.Fprintln(os.Stderr, result.String()) //If Error != nil then the String call prints the error
		} else if len(result.IP) == 1 {
				fmt.Println(result.Json())
		} else {
			fmt.Println(result.Json())
		}
	}
}

func lookupDomain(domain string, results chan dnsResult)  {
	defer wg.Done()
	defer wgInner.Done()
	ips, err := net.LookupIP(domain)
	if err != nil {
		results <- dnsResult{
			Domain: domain,
			IP:     nil,
			Error:  err,
		}
	} else {
		var temp []net.IP
		for _, ip := range ips {
			temp = append(temp, net.ParseIP(ip.String()))
		}
		results <- dnsResult{
			Domain: domain,
			IP: temp,
			Error:  nil,
		}
	}
}
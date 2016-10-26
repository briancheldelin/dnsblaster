package main

import (
	"fmt"
	"os"
	"time"

	"github.com/miekg/dns"
)

func main() {
	target := os.Args[1]
	service := os.Args[2]
	port := os.Args[3]

	target = target + ":" + port
	fmt.Printf("Target: %v\n", target)

	blasterCount := 100
	var total int
	var totalError int

	totals := make(chan int, blasterCount)
	totalsBad := make(chan int, blasterCount)

	// Blaster Loop
	for i := 0; i < blasterCount; i++ {

		go func() {
			i := 0
			b := 0
			start := time.Now()
			for {
				m := new(dns.Msg)
				m.SetQuestion(dns.Fqdn(service), dns.TypeA)
				in, err := dns.Exchange(m, target)
				if err != nil {
					b++
				} else {
					if in != nil {
						i++
					}
				}
				if time.Since(start).Seconds() >= 120 {
					break
				}
			}
			fmt.Printf("Count: %v\n", i)
			totals <- i
			totalsBad <- b

		}() //Function end

	}

	for i := 0; i < blasterCount; i++ {
		p := <-totals
		total = p + total
	}

	for i := 0; i < blasterCount; i++ {
		p := <-totalsBad
		totalError = p + totalError
	}
	fmt.Printf("Total good: %d\n", total)
	fmt.Printf("Total bad: %d\n", totalError)
	fmt.Printf("Request Per Min: %d\n", total/2)
	fmt.Printf("Request Per Sec: %d\n", total/120)

}

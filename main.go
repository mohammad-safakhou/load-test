package main

import (
	"flag"
	"fmt"
	"load-test/engine"
	"time"
)

func main() {
	start := time.Now()
	numberOfAccounts := flag.Int("num", 1000000, "Number of accounts to create")
	timeout := flag.Duration("timeout", 5*time.Second, "Number of accounts to create")
	flag.Parse()
	fmt.Printf("Number of accounts to create: %d\n", *numberOfAccounts)
	engine.Start(*numberOfAccounts, *timeout)
	fmt.Printf("Time elapsed: %v\n", time.Since(start))
}

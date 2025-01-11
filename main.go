package main

import (
	"flag"
	"fmt"
	"load-test/engine"
	"load-test/monitoring"
)

func main() {
	numberOfAccounts := flag.Int("num", 1000000, "Number of accounts to create")
	flag.Parse()
	fmt.Printf("Number of accounts to create: %d\n", *numberOfAccounts)
	//monitoring.Init(*numberOfAccounts)
	engine.Start(*numberOfAccounts)
	close(monitoring.Monitoring)
}

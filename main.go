package main

import (
	"flag"
	"fmt"
	"load-test/engine"
)

func main() {
	numberOfAccounts := flag.Int("num", 1000000, "Number of accounts to create")
	flag.Parse()
	fmt.Printf("Number of accounts to create: %d\n", *numberOfAccounts)
	engine.Start(*numberOfAccounts)
}

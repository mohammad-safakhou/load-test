package main

import (
	"fmt"
	"load-test/diameter"
	"load-test/engine"
	"time"
)

func main() {
	diameter.CCAs = []string{}
	numberOfAccounts := 10000
	engine.Start(numberOfAccounts)
	for {
		time.Sleep(1 * time.Second)
		if len(diameter.CCAs) == numberOfAccounts {
			break
		} else {
			fmt.Printf("%d CCA(s) fetched.\n", len(diameter.CCAs))
		}
	}
}

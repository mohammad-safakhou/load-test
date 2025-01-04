package main

import (
	"fmt"
	"load-test/diameter"
	"load-test/engine"
	"load-test/monitoring"
	"time"
)

func main() {
	diameter.CCAs = []string{}
	numberOfAccounts := 10000
	monitoring.Init(numberOfAccounts)
	engine.Start(numberOfAccounts)
	for {
		time.Sleep(1 * time.Second)
		if len(diameter.CCAs) == numberOfAccounts {
			break
		} else {
			fmt.Printf("%d CCA(s) fetched.\n", len(diameter.CCAs))
		}
	}
	fmt.Printf("%d CCA(s) fetched.\n", len(diameter.CCAs))
	close(monitoring.Monitoring)
}

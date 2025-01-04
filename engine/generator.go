package engine

import (
	"fmt"
	"github.com/MHG14/go-diameter/v4/diam"
	"load-test/diameter"
	"load-test/pipeline"
	"math"
	"strconv"
	"sync"
	"time"
)

var conn diam.Conn
var client diameter.Client

const batchSize = 1

const updateIterations = 2
const sleepTimes = 1 * time.Second

func worker(task chan string, wg *sync.WaitGroup, numberOfAccounts int) {
	defer wg.Done()
	for id := range task {
		pipeline.NewAccount(updateIterations, numberOfAccounts, sleepTimes, client, id).Run()
		fmt.Printf("%s is Done\n", id)
		return
	}
}

func Start(numberOfAccounts int) {
	var err error
	conn, err = diameter.NewConnection()
	if err != nil {
		panic(err)
	}
	client = diameter.NewDiameterClient(conn)

	tasks := make(chan string, numberOfAccounts)
	wg := new(sync.WaitGroup)
	wg.Add(numberOfAccounts)
	for i := 0; i < numberOfAccounts; i++ {
		go worker(tasks, wg, numberOfAccounts)
	}

	fmt.Println("Workers are all up and running")

	var accountIDs []string
	for i := 1; i <= numberOfAccounts; i++ {
		id := "00" + strconv.Itoa(i)
		accountIDs = append(accountIDs, id)
	}

	accountsMap := make(map[int][]string)

	numOfWorkers := int(math.Ceil(float64(numberOfAccounts) / batchSize))

	for i := 0; i < numOfWorkers; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > numberOfAccounts {
			end = numberOfAccounts
		}
		accountsMap[i] = accountIDs[start:end]
	}

	pusherWg := new(sync.WaitGroup)
	pusherWg.Add(numOfWorkers)
	for i := 0; i < numOfWorkers; i++ {
		go pushWorker(tasks, accountsMap[i], pusherWg)
	}

	pusherWg.Wait()

	close(tasks)
	wg.Wait()
}

func pushWorker(tasks chan string, accountIDs []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, id := range accountIDs {
		tasks <- id
	}
}

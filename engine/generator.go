package engine

import (
	"fmt"
	"load-test/diameter"
	"load-test/models"
	"load-test/pipeline"
	"math"
	"sync"
	"time"
)

const batchSize = 1

const updateIterations = 2
const sleepTimes = 1 * time.Second

func worker(task chan models.AccountID, wg *sync.WaitGroup, numberOfAccounts int, client diameter.Client) {
	defer wg.Done()
	for id := range task {
		pipeline.NewAccount(updateIterations, numberOfAccounts, sleepTimes, client, id).Run()
		fmt.Printf("%s is Done\n", id)
	}
}

func Start(numberOfAccounts int, timeout time.Duration) {
	var err error
	hopIDs := new(sync.Map)
	conn, err := diameter.NewConnection(hopIDs)
	if err != nil {
		panic(err)
	}
	client := diameter.NewDiameterClient(conn, hopIDs, timeout)

	tasks := make(chan models.AccountID, numberOfAccounts)
	wg := new(sync.WaitGroup)
	wg.Add(numberOfAccounts)
	for i := 0; i < numberOfAccounts; i++ {
		go worker(tasks, wg, numberOfAccounts, client)
	}

	fmt.Println("Workers are all up and running")

	var accountIDs []models.AccountID
	for i := 1; i <= numberOfAccounts; i++ {
		id := models.NewAccountID(i)
		accountIDs = append(accountIDs, id)
	}

	accountsMap := make(map[int][]models.AccountID)

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

func pushWorker(tasks chan models.AccountID, accountIDs []models.AccountID, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, id := range accountIDs {
		tasks <- id
	}
}

package pipeline

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"load-test/diameter"
	"strconv"
	"sync"
	"time"
)

type Launcher interface {
	Run()
}

type account struct {
	updateIteration     int
	sleepTimes          time.Duration
	client              diameter.Client
	accountID           string
	otherID             string
	sessionData         string
	sessionVoiceCalling string
	sessionVoiceCalled  string
	sessionVideoCalling string
	sessionVideoCalled  string
}

func NewAccount(
	updateIteration int,
	numberOfAccounts int,
	sleepTimes time.Duration,
	client diameter.Client,
	accountID string,
) Launcher {
	ai, _ := strconv.Atoi(accountID[2:])
	temp := numberOfAccounts - ai + 1
	if temp == ai {
		temp += 1
	}
	return &account{
		updateIteration: updateIteration,
		sleepTimes:      sleepTimes,
		client:          client,
		accountID:       accountID,
		otherID:         "00" + strconv.Itoa(temp),
	}
}

func (m *account) Run() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.sessionData = fmt.Sprintf("%s:10:%s", m.accountID, uuid.New().String())
		err := m.client.InitData(m.accountID, m.sessionData)
		if err != nil {
			log.Errorf("init data err: %v", err)
		}

		//time.Sleep(10 * time.Second)

		for i := 0; i < m.updateIteration; i++ {
			time.Sleep(m.sleepTimes)
			err = m.client.UpdateData(m.accountID, m.sessionData)
			if err != nil {
				log.Errorf("update data err: %v", err)
			}
		}

		err = m.client.TerminateData(m.accountID, m.sessionData)
		if err != nil {
			log.Errorf("terminate data err: %v", err)
		}
	}()
	//
	//go func() {
	//	defer wg.Done()
	//	m.sessionVoiceCalling = fmt.Sprintf("%s:20:%s", m.accountID, uuid.New().String())
	//	err := m.client.InitVoiceCalling(m.accountID, m.otherID, m.sessionVoiceCalling)
	//	if err != nil {
	//		log.Errorf("init voice calling err: %v", err)
	//	}
	//	m.sessionVoiceCalled = fmt.Sprintf("%s:21:%s", m.accountID, uuid.New().String())
	//	err = m.client.InitVoiceCalled(m.accountID, m.otherID, m.sessionVoiceCalled)
	//	if err != nil {
	//		log.Errorf("init voice called err: %v", err)
	//	}
	//
	//	innerWG := new(sync.WaitGroup)
	//	innerWG.Add(2)
	//	go func() {
	//		defer innerWG.Done()
	//		for i := 0; i < m.updateIteration; i++ {
	//			time.Sleep(m.sleepTimes)
	//			err = m.client.UpdateVoiceCalling(m.accountID, m.otherID, m.sessionVoiceCalling)
	//			if err != nil {
	//				log.Errorf("update voice calling err: %v", err)
	//			}
	//		}
	//	}()
	//	go func() {
	//		defer innerWG.Done()
	//		for i := 0; i < m.updateIteration; i++ {
	//			time.Sleep(m.sleepTimes)
	//			err = m.client.UpdateVoiceCalled(m.accountID, m.otherID, m.sessionVoiceCalled)
	//			if err != nil {
	//				log.Errorf("update voice called err: %v", err)
	//			}
	//		}
	//	}()
	//	innerWG.Wait()
	//
	//	err = m.client.TerminateVoiceCalling(m.accountID, m.otherID, m.sessionVoiceCalling)
	//	if err != nil {
	//		log.Errorf("terminate voice calling err: %v", err)
	//	}
	//	err = m.client.TerminateVoiceCalled(m.accountID, m.otherID, m.sessionVoiceCalled)
	//	if err != nil {
	//		log.Errorf("terminate voice called err: %v", err)
	//	}
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//	m.sessionVideoCalling = fmt.Sprintf("%s:20:%s", m.accountID, uuid.New().String())
	//	err := m.client.InitVideoCalling(m.accountID, m.otherID, m.sessionVideoCalling)
	//	if err != nil {
	//		log.Errorf("init video calling err: %v", err)
	//	}
	//
	//	for i := 0; i < m.updateIteration; i++ {
	//		time.Sleep(m.sleepTimes)
	//		err = m.client.UpdateVideoCalling(m.accountID, m.otherID, m.sessionVideoCalling)
	//		if err != nil {
	//			log.Errorf("update video calling err: %v", err)
	//		}
	//	}
	//
	//	err = m.client.TerminateVideoCalling(m.accountID, m.otherID, m.sessionVideoCalling)
	//	if err != nil {
	//		log.Errorf("terminate video calling err: %v", err)
	//	}
	//}()

	wg.Wait()
}

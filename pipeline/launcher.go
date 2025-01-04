package pipeline

import (
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"load-test/diameter"
	"sync"
)

type meta struct {
	client              diameter.Client
	AccountID           string
	SessionData         string
	SessionVoiceCalling string
	SessionVoiceCalled  string
	SessionVideoCalling string
	SessionVideoCalled  string
}

func (m *meta) Run() error {
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go func() {
		defer wg.Done()
		m.SessionData = fmt.Sprintf("%s:10:%s", m.AccountID, uuid.New().String())
		err := m.client.InitData(m.AccountID, m.SessionData)
		if err != nil {
			log.Errorf("init data err: %v", err)
		}
		m.SessionVideoCalling = fmt.Sprintf("%s:10:%s", m.AccountID, uuid.New().String())
		err := m.client.InitData(m.AccountID, m.SessionData)
		if err != nil {
			log.Errorf("init data err: %v", err)
		}
		m.SessionData = fmt.Sprintf("%s:10:%s", m.AccountID, uuid.New().String())
		err := m.client.InitData(m.AccountID, m.SessionData)
		if err != nil {
			log.Errorf("init data err: %v", err)
		}
	}()
}

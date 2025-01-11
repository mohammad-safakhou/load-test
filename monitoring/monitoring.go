package monitoring

import "sync"

var monit map[string]Monit

type Monit struct {
	AccountID   string
	IsInit      bool
	IsTerminate bool
	UpdateCount int
	MaxUpdate   int

	mu sync.Mutex
}

type State struct {
	AccountID string
	Type      string
}

type StateType string

const (
	StateTypeInit      = "init"
	StateTypeTerminate = "terminate"
	StateTypeUpdate    = "update"
)

var Monitoring chan State

func Init(numberOfAccounts int) func() {
	Monitoring = make(chan State, numberOfAccounts)
	for i := 0; i < 10; i++ {
		go monitoringWorker()
	}
	return closer
}

func closer() {
	close(Monitoring)
}

func monitoringWorker() {
	for state := range Monitoring {
		handler(state)
	}
}

func handler(s State) {
	if m, ok := monit[s.AccountID]; ok {
		m.mu.Lock()
		defer m.mu.Unlock()

		if s.Type == StateTypeUpdate {
			m.UpdateCount += 1
		} else if s.Type == StateTypeTerminate {
			m.IsTerminate = true
		} else if s.Type == StateTypeInit {
			panic("why we are here???")
		} else {
			panic("woooooow we should not be here")
		}
	} else {
		monit[s.AccountID] = Monit{
			AccountID:   s.AccountID,
			IsInit:      true,
			IsTerminate: false,
			UpdateCount: 0,
			MaxUpdate:   0,
			mu:          sync.Mutex{},
		}
	}
}

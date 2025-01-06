package diameter

import (
	"github.com/MHG14/go-diameter/v4/diam"
	"time"
)

type MockDiameterClient struct {
	conn diam.Conn
	//mux  *sm.StateMachine
}

func (d *MockDiameterClient) InitData(accountID string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "initial")
	}()
	return nil
}

func (d *MockDiameterClient) UpdateData(accountID string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "update")
	}()
	return nil
}

func (d *MockDiameterClient) TerminateData(accountID string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "terminate")
	}()
	return nil
}

func (d *MockDiameterClient) InitVideoCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "initial")
	}()
	return nil
}

func (d *MockDiameterClient) UpdateVideoCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "update")
	}()
	return nil
}

func (d *MockDiameterClient) TerminateVideoCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "terminate")
	}()
	return nil
}

func (d *MockDiameterClient) InitVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "initial")
	}()
	return nil
}

func (d *MockDiameterClient) UpdateVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "update")
	}()
	return nil
}

func (d *MockDiameterClient) TerminateVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "terminate")
	}()
	return nil
}

func (d *MockDiameterClient) InitVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "initial")
	}()
	return nil
}

func (d *MockDiameterClient) UpdateVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "update")
	}()
	return nil
}

func (d *MockDiameterClient) TerminateVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	time.Sleep(1 * time.Second)
	go func() {
		CCAs = append(CCAs, "terminate")
	}()
	return nil
}

func NewMockDiameterClient(conn diam.Conn) Client {
	return &MockDiameterClient{
		conn: conn,
	}
}

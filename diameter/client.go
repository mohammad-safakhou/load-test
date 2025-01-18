package diameter

import (
	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/pkg/errors"
	"load-test/models"
	"sync"
	"time"
)

type Client interface {
	InitData(accountID models.AccountID, sessionID string) error
	UpdateData(accountID models.AccountID, sessionID string) error
	TerminateData(accountID models.AccountID, sessionID string) error

	InitVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error
	UpdateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error
	TerminateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error

	InitVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error
	UpdateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error
	TerminateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error
	InitVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error
	UpdateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error
	TerminateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error
}

type DiameterClient struct {
	timeout time.Duration
	conn    diam.Conn

	hopIDs *sync.Map
	//mux  *sm.StateMachine
}

func (d *DiameterClient) Send(message *diam.Message, accountID models.AccountID) error {
	hopID := message.Header.HopByHopID
	ch := make(chan *diam.Message)

	d.hopIDs.Store(hopID, ch)

	conn, err := NewConnection(d.hopIDs)
	if err != nil {
		panic(errors.Wrap(err, "unable to connect to diameter in send function"))
	}
	defer conn.Close()
	_, err = message.WriteTo(conn)
	if err != nil {
		return err
	}

	timeout := time.After(d.timeout)

	// Wait for Response
	select {
	case resp := <-ch:
		_ = resp
		d.hopIDs.Delete(hopID)

		return nil
	case <-timeout:
		d.hopIDs.Delete(hopID)
		//return errors.New(fmt.Sprintf("Timeout happened on accountID: %s", accountID.String()))
		return nil
	}
}

func (d *DiameterClient) InitData(accountID models.AccountID, sessionID string) error {
	return d.Send(BuildDataInitSessionCCR(sessionID, accountID.String()), accountID)
}

func (d *DiameterClient) UpdateData(accountID models.AccountID, sessionID string) error {
	return d.Send(BuildDataUpdateSessionCCR(sessionID, accountID.String()), accountID)
}

func (d *DiameterClient) TerminateData(accountID models.AccountID, sessionID string) error {
	return d.Send(BuildDataTerminateSessionCCR(sessionID, accountID.String()), accountID)
}

func (d *DiameterClient) InitVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVideoCallingInitSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) UpdateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVideoCallingUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) TerminateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVideoCallingTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) InitVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCallingInitSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) UpdateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCallingUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) TerminateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCallingTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) InitVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCalledInitSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) UpdateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCalledUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func (d *DiameterClient) TerminateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	return d.Send(BuildVoiceCalledTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String()), accountID0)
}

func NewDiameterClient(conn diam.Conn, hopIDs *sync.Map, timeout time.Duration) Client {
	return &DiameterClient{
		timeout: timeout,
		conn:    conn,
		hopIDs:  hopIDs,
	}
}

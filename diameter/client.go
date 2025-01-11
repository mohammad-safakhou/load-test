package diameter

import (
	"github.com/MHG14/go-diameter/v4/diam"
	"load-test/models"
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
	conn diam.Conn
	//mux  *sm.StateMachine
}

func (d *DiameterClient) InitData(accountID models.AccountID, sessionID string) error {
	message := BuildDataInitSessionCCR(sessionID, accountID.String())
	_, err := message.WriteToStream(d.conn, uint(accountID.ID()))
	return err
}

func (d *DiameterClient) UpdateData(accountID models.AccountID, sessionID string) error {
	message := BuildDataUpdateSessionCCR(sessionID, accountID.String())
	_, err := message.WriteToStream(d.conn, uint(accountID.ID()))
	return err
}

func (d *DiameterClient) TerminateData(accountID models.AccountID, sessionID string) error {
	message := BuildDataTerminateSessionCCR(sessionID, accountID.String())
	_, err := message.WriteToStream(d.conn, uint(accountID.ID()))
	return err
}

func (d *DiameterClient) InitVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVideoCallingInitSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) UpdateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVideoCallingUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) TerminateVideoCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVideoCallingTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) InitVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCallingInitSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) UpdateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCallingUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) TerminateVoiceCalling(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCallingTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) InitVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCalledInitSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) UpdateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCalledUpdateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func (d *DiameterClient) TerminateVoiceCalled(accountID0, accountID1 models.AccountID, sessionID string) error {
	message := BuildVoiceCalledTerminateSessionCCR(sessionID, accountID0.String(), accountID1.String())
	_, err := message.WriteToStream(d.conn, uint(accountID0.ID()))
	return err
}

func NewDiameterClient(conn diam.Conn) Client {
	return &DiameterClient{
		conn: conn,
	}
}

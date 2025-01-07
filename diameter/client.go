package diameter

import "github.com/MHG14/go-diameter/v4/diam"

type Client interface {
	InitData(accountID string, sessionID string) error
	UpdateData(accountID string, sessionID string) error
	TerminateData(accountID string, sessionID string) error

	InitVideoCalling(accountID0, accountID1 string, sessionID string) error
	UpdateVideoCalling(accountID0, accountID1 string, sessionID string) error
	TerminateVideoCalling(accountID0, accountID1 string, sessionID string) error

	InitVoiceCalling(accountID0, accountID1 string, sessionID string) error
	UpdateVoiceCalling(accountID0, accountID1 string, sessionID string) error
	TerminateVoiceCalling(accountID0, accountID1 string, sessionID string) error
	InitVoiceCalled(accountID0, accountID1 string, sessionID string) error
	UpdateVoiceCalled(accountID0, accountID1 string, sessionID string) error
	TerminateVoiceCalled(accountID0, accountID1 string, sessionID string) error
}

type DiameterClient struct {
	conn diam.Conn
	//mux  *sm.StateMachine
}

func (d *DiameterClient) InitData(accountID string, sessionID string) error {
	message := BuildDataInitSessionCCR(sessionID, accountID)
	cnn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	_, err = message.WriteTo(cnn)
	return err
}

func (d *DiameterClient) UpdateData(accountID string, sessionID string) error {
	message := BuildDataUpdateSessionCCR(sessionID, accountID)
	cnn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	_, err = message.WriteTo(cnn)
	return err
}

func (d *DiameterClient) TerminateData(accountID string, sessionID string) error {
	message := BuildDataTerminateSessionCCR(sessionID, accountID)
	cnn, err := NewConnection()
	if err != nil {
		panic(err)
	}
	_, err = message.WriteTo(cnn)
	return err
}

func (d *DiameterClient) InitVideoCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVideoCallingInitSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) UpdateVideoCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVideoCallingUpdateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) TerminateVideoCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVideoCallingTerminateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) InitVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCallingInitSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) UpdateVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCallingUpdateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) TerminateVoiceCalling(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCallingTerminateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) InitVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCalledInitSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) UpdateVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCalledUpdateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func (d *DiameterClient) TerminateVoiceCalled(accountID0, accountID1 string, sessionID string) error {
	message := BuildVoiceCalledTerminateSessionCCR(sessionID, accountID0, accountID1)
	_, err := message.WriteTo(d.conn)
	return err
}

func NewDiameterClient(conn diam.Conn) Client {
	return &DiameterClient{
		conn: conn,
	}
}

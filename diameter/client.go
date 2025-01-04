package diameter

type Client interface {
	InitData(accountID string, sessionID string) error
	UpdateData(accountID string, sessionID string) error
	TerminateData(accountID string, sessionID string) error

	InitVideoCalling(accountID string, sessionID string) error
	UpdateVideoCalling(accountID string, sessionID string) error
	TerminateVideoCalling(accountID string, sessionID string) error
	InitVideoCalled(accountID string, sessionID string) error
	UpdateVideoCalled(accountID string, sessionID string) error
	TerminateVideoCalled(accountID string, sessionID string) error

	InitVoiceCalling(accountID string, sessionID string) error
	UpdateVoiceCalling(accountID string, sessionID string) error
	TerminateVoiceCalling(accountID string, sessionID string) error
	InitVoiceCalled(accountID string, sessionID string) error
	UpdateVoiceCalled(accountID string, sessionID string) error
	TerminateVoiceCalled(accountID string, sessionID string) error
}

type DiameterClient struct {
	//conn diam.Conn
	//mux  *sm.StateMachine
}

//func (dc *DiameterClient) SendSNR(message *service_models.SNRMessage) error {
//	fmt.Printf("SN REQUEST TYPE IS %+v\n", message)
//	// Get this client's metadata from the connection object,
//	// which is set by the state machine after the handshake.
//	// It contains the peer's Origin-Host and Realm from the
//	// CER/CEA handshake. We use it to populate the AVPs below.
//
//	//meta, ok := smpeer.FromContext(c.Context()
//	//if !ok {
//	//  return errors.New("peer metadata unavailable")
//	//}
//
//	channel := utils.SetSessionChannel(message.SessionId)
//	defer close(channel)
//
//	m := diam.NewRequest(8388636, 16777302, dc.conn.Dictionary()) // SLR ID,Diameter Sy ID
//	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(message.SessionId))
//	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.UTF8String(config.AppConfig.Agents.Diameter.Client.OriginHost))
//	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.UTF8String(config.AppConfig.Agents.Diameter.Client.OriginRealm))
//	m.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.UTF8String(config.AppConfig.Agents.RemoteDiameter.OriginHost))
//	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.UTF8String(config.AppConfig.Agents.RemoteDiameter.OriginRealm))
//	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(16777302))
//	m.NewAVP(avp.SNRequestType, avp.Mbit, 0, datatype.Unsigned32(message.SNRequestType))
//	for _, pcsr := range message.PolicyCounterStatusReport {
//		m.NewAVP(avp.PolicyCounterStatusReport, avp.Mbit, 0, &diam.GroupedAVP{
//			AVP: []*diam.AVP{
//				diam.NewAVP(avp.PolicyCounterIdentifier, avp.Mbit, 0, datatype.UTF8String(pcsr.PolicyCounterIdentifier)),
//				diam.NewAVP(avp.PolicyCounterStatus, avp.Mbit, 0, datatype.UTF8String(pcsr.PolicyCounterStatus)),
//				diam.NewAVP(avp.PendingPolicyCounterInformation, avp.Mbit, 0, &diam.GroupedAVP{
//					AVP: []*diam.AVP{
//						diam.NewAVP(avp.PolicyCounterStatus, avp.Mbit, 0, datatype.UTF8String(pcsr.PendingPolicyCounterInformation.PolicyCounterStatus)),
//						diam.NewAVP(avp.PendingPolicyCounterChangeTime, avp.Mbit, 0, datatype.Time(pcsr.PendingPolicyCounterInformation.PendingPolicyCounterChangeTime)),
//					},
//				}),
//			},
//		})
//	}
//	log.Printf("Sending SLR to %s\n%s", dc.conn.RemoteAddr(), m)
//	_, err := m.WriteTo(dc.conn)
//
//	if err == nil {
//		fmt.Println("successfully wrote SNR to connection")
//	}
//
//	//done := make(chan diam.Message, 1000)
//	//
//	////dc.mux.Handle("SNA", handleSNA(done))
//	//diam.Handle("ALL", handleSNA1(done))
//	//diam.Handle("SNA", handleSNA2(done))
//	//diam.Handle("Spending-Status-Notification-Answer", handleSNA3(done))
//	//
//
//	//_ = <-done
//
//	deadline := time.Tick(time.Second * 20)
//
//	for {
//		select {
//		case sna := <-channel:
//			fmt.Println("incoming sna from pcrf issssssssssss", sna)
//			utils.RemoveSessionChannel(message.SessionId)
//			return nil
//		case <-deadline:
//			return errors.New("no sna received")
//		}
//	}
//
//	//return err
//}

//func handleSNA(done chan diam.Message) diam.HandlerFunc {
//	return func(c diam.Conn, m *diam.Message) {
//		fmt.Printf("THE CONNECTION OBJECT FROM HANDLESNA IS %+v\n", c)
//		log.Printf("Received fucking SNA from %s\n%s", c.RemoteAddr(), m)
//		done <- *m
//	}
//}

//func NewDiameterClient(conn diam.Conn) *DiameterClient {
//	return &DiameterClient{
//		conn: conn,
//	}
//}

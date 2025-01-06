package diameter

import (
	"flag"
	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
	"github.com/MHG14/go-diameter/v4/diam/sm"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func NewConnection() (diam.Conn, error) {
	addr := flag.String("addr", "192.168.20.244:3868", "address in form of ip:port to connect to")
	ssl := flag.Bool("ssl", false, "connect to server using tls")
	host := flag.String("diam_host", "client", "diameter identity host")
	realm := flag.String("diam_realm", "go-diameter", "diameter identity realm")
	certFile := flag.String("cert_file", "", "tls client certificate file (optional)")
	keyFile := flag.String("key_file", "", "tls client key file (optional)")
	networkType := flag.String("network_type", "tcp", "protocol type tcp/sctp")

	cfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity(*host),
		OriginRealm:      datatype.DiameterIdentity(*realm),
		VendorID:         0,
		ProductName:      "go-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(net.ParseIP("127.0.0.1")),
		},
	}

	// Create the state machine (it's a diam.ServeMux) and client.
	mux := sm.New(cfg)

	mux.Handle("CCA", handleCCA())

	cli := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     false,
		WatchdogInterval:   5 * time.Second,
		WatchdogStream:     0,
		SupportedVendorID:  nil,
		AcctApplicationID:  nil,
		//AcctApplicationID: []*diam.AVP{
		//	// Advertise that we want support accounting application with id 999
		//	diam.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		//},
		AuthApplicationID: []*diam.AVP{
			// Advertise support for credit control application
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)), // RFC 4006
		},
		VendorSpecificApplicationID: nil,
	}
	return dial(cli, *addr, *certFile, *keyFile, *ssl, *networkType)
}

func dial(cli *sm.Client, addr, cert, key string, ssl bool, networkType string) (diam.Conn, error) {
	if ssl {
		return cli.DialNetworkTLS(networkType, addr, cert, key, nil)
	}
	return cli.DialNetwork(networkType, addr)
}

type CCAMessage struct {
	RequestType datatype.Unsigned32 `avp:"CC-Request-Type"`
}

var CCAs []string

func handleCCA() diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		message := CCAMessage{}
		err := m.Unmarshal(&message)
		if err != nil {
			log.Println(err)
			return
		}
		//log.Printf("Received CCA from %s ------- %s", c.RemoteAddr(), message.RequestType.String())
		if message.RequestType == datatype.Unsigned32(3) {
			CCAs = append(CCAs, message.RequestType.String())
		}
	}
}
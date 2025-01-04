package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

const Abbas = 10415

// Example structures to hold config/user data:
type Config struct {
	MNC                string
	MCC                string
	BitrateUL          uint32
	BitrateDL          uint32
	PreMiddleFix       string
	Prefix             string
	TimeZone           string
	UserLocationInfo   string
	DownloadedBytes    uint32
	UploadedBytes      uint32
	VideoRequestedTime uint32
}

// BuildInitSessionCCR builds a CCR message (Credit-Control Request)
func BuildDataInitSessionCCR(
	sessionID string,
	phoneNumber string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// 1) Create a new CCR message: Command-Code=272 (Credit-Control), App-ID=4 (Gx).
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 2) Add standard AVPs

	// Session-Id
	sessionIDVal := fmt.Sprintf(
		"smf.epc.mnc%s.mcc%s.3gppnetwork.org;%s;%s",
		cfg.MNC, cfg.MCC, sessionID, sessionID,
	)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionIDVal))

	// Origin-Host
	originHostVal := fmt.Sprintf(
		"smf.epc.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(originHostVal))

	// Origin-Realm
	originRealmVal := fmt.Sprintf(
		"epc.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(originRealmVal))

	// Destination-Realm
	destRealmVal := fmt.Sprintf(
		"epc.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.DiameterIdentity(destRealmVal))

	// Auth-Application-Id (4 for Gx)
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))

	// Service-Context-Id: 461 is typical AVP code for Service-Context-Id
	m.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("32251@3gpp.org"))

	// CC-Request-Type: 416
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1)) // 1=INITIAL_REQUEST

	// CC-Request-Number: 415
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))

	// Event-Timestamp: 55
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 3) Add Subscription-Id AVPs (Grouped)

	// Subscription-Id #1
	subIDGrouped1 := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)), // Type=0
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(phoneNumber)),
		},
	}
	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, subIDGrouped1)

	// Subscription-Id #2
	subIDGrouped2 := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(1)), // Type=1
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(
				cfg.PreMiddleFix+phoneNumber,
			)),
		},
	}
	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, subIDGrouped2)

	// Subscription-Id #3 (repeated type=0 from your example)
	subIDGrouped3 := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(phoneNumber)),
		},
	}
	m.NewAVP(avp.SubscriptionID, avp.Mbit, 0, subIDGrouped3)

	// Requested-Action: 436 (0=CHECK_BALANCE or DIRECT_DEBIT, depends on spec)
	m.NewAVP(avp.RequestedAction, avp.Mbit, 0, datatype.Enumerated(0))

	// AoC-Request-Type (Vendor-Specific, example: code=812, vendor=10415)
	m.NewAVP(
		avp.AoCRequestType,
		avp.Mbit|avp.Vbit,
		Abbas,
		datatype.Enumerated(1), // e.g., 1 = AoC request type
	)

	// Multiple-Services-Indicator: 455 (0 = MULTIPLE_SERVICES_NOT_SUPPORTED)
	m.NewAVP(avp.MultipleServicesIndicator, avp.Mbit, 0, datatype.Enumerated(0))

	// Multiple-Services-Credit-Control (Grouped): 456
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Requested-Service-Unit (437)
			diam.NewAVP(
				avp.RequestedServiceUnit,
				avp.Mbit,
				0,
				datatype.UTF8String(""), // from your snippet (empty string)
			),
			// Used-Service-Unit (446) -> grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(0)),
						diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(0)),
						diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(0)),
					},
				},
			),
			// QoS-Information (3GPP vendor-specific)
			diam.NewAVP(
				avp.QoSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// QoS-Class-Identifier (QCI=5), example code=1028
						diam.NewAVP(1028, avp.Mbit|avp.Vbit, Abbas, datatype.Enumerated(5)),
						// Allocation-Retention-Priority (code=1034)
						diam.NewAVP(1034, avp.Vbit, Abbas, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								// Priority-Level=1046
								diam.NewAVP(1046, avp.Vbit, Abbas, datatype.Enumerated(1)),
								// Preemption-Capability=1047
								diam.NewAVP(1047, avp.Vbit, Abbas, datatype.Enumerated(1)),
								// Preemption-Vulnerability=1048
								diam.NewAVP(1048, avp.Vbit, Abbas, datatype.Enumerated(1)),
							},
						}),
						// APN-Aggregate-Max-Bitrate-UL=1041
						diam.NewAVP(1041, avp.Vbit, Abbas, datatype.Unsigned32(cfg.BitrateUL)),
						// APN-Aggregate-Max-Bitrate-DL=1042
						diam.NewAVP(1040, avp.Vbit, Abbas, datatype.Unsigned32(cfg.BitrateDL)),
					},
				},
			),
			// TGPP-RAT-Type (code=1032): "06" as OctetString
			diam.NewAVP(
				avp.TGPPRATType,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.OctetString("06"),
			),
		},
	}
	m.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, msccGrouped)

	// Service-Information (vendor-specific)
	psInfoGroupedAVP := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// TGPP-Charging-ID (example code=2)
			diam.NewAVP(2, avp.Mbit|avp.Vbit, Abbas, datatype.OctetString("0000000a")),
			// TGPP-PDP-Type (example code=3)
			diam.NewAVP(3, avp.Mbit|avp.Vbit, Abbas, datatype.Enumerated(0)),
			// PDP-Address (example code=1227)
			diam.NewAVP(1227, avp.Mbit|avp.Vbit, Abbas, datatype.Address("10.46.0.2")),
			// SGSN-Address (example code=1228)
			diam.NewAVP(1228, avp.Mbit|avp.Vbit, Abbas, datatype.Address("127.0.0.3")),
			// GGSN-Address (example code=1229) - repeated
			diam.NewAVP(1229, avp.Mbit|avp.Vbit, Abbas, datatype.Address("127.0.0.4")),
			diam.NewAVP(1229, avp.Mbit|avp.Vbit, Abbas, datatype.Address("127.0.0.4")),
			// Called-Station-Id (30)
			diam.NewAVP(avp.CalledStationID, avp.Mbit, 0, datatype.UTF8String("internet")),
			// TGPP-Selection-Mode (example code=4)
			diam.NewAVP(4, avp.Mbit|avp.Vbit, Abbas, datatype.UTF8String("0")),
			// TGPP-SGSN-MCCMNC (example code=18)
			diam.NewAVP(18, avp.Mbit|avp.Vbit, Abbas, datatype.UTF8String(cfg.Prefix)),
			// TGPP-NSAPI (example code=5)
			diam.NewAVP(5, avp.Mbit|avp.Vbit, Abbas, datatype.OctetString("\\005")),
			// TGPP-MS-TimeZone (example code=23)
			diam.NewAVP(23, avp.Mbit|avp.Vbit, Abbas, datatype.OctetString(cfg.TimeZone)),
			// TGPP-User-Location-Info (example code=22)
			diam.NewAVP(22, avp.Mbit|avp.Vbit, Abbas, datatype.UTF8String(cfg.UserLocationInfo)),
			// User-Equipment-Info (standard code=458 in RFC, but 3GPP extends it)
			diam.NewAVP(
				avp.UserEquipmentInfo,
				avp.Mbit,
				0, // no vendor ID for the standard AVP
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(
							avp.UserEquipmentInfoType,
							avp.Mbit,
							0,
							datatype.Enumerated(0),
						), // 0 = IMEI
						diam.NewAVP(
							avp.UserEquipmentInfoValue,
							avp.Mbit,
							0,
							datatype.OctetString("doops"),
						),
					},
				},
			),
		},
	}

	serviceInfoGroupedAVP := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.PSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				psInfoGroupedAVP,
			),
		},
	}

	m.NewAVP(
		avp.ServiceInformation,
		avp.Mbit|avp.Vbit,
		Abbas,
		serviceInfoGroupedAVP,
	)

	// Return the fully built CCR message
	return m
}

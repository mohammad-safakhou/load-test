package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildVoiceCallingTerminateSessionCCR parallels your Node.js "voiceCallingTerminate" code in Go.
func BuildVoiceCallingTerminateSessionCCR(
	sessionID string,
	phoneNumberCalling string,
	phoneNumberCalled string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// Create a CCR message: Command-Code=272 (Credit-Control), Application-ID=4
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 1) Basic / standard AVPs

	// Session-Id
	sessionIDVal := fmt.Sprintf(
		"smf.epc.mnc%s.mcc%s.3gppnetwork.org;%s;%s",
		cfg.MNC, cfg.MCC, sessionID, sessionID,
	)
	m.NewAVP(
		avp.SessionID,
		avp.Mbit,
		0,
		datatype.UTF8String(sessionIDVal),
	)

	// Origin-Host
	originHostVal := fmt.Sprintf(
		"scscf.ims.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(
		avp.OriginHost,
		avp.Mbit,
		0,
		datatype.UTF8String(originHostVal), // or DiameterIdentity
	)

	// Origin-Realm
	originRealmVal := fmt.Sprintf(
		"ims.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(
		avp.OriginRealm,
		avp.Mbit,
		0,
		datatype.UTF8String(originRealmVal), // or DiameterIdentity
	)

	// Destination-Realm
	destRealmVal := fmt.Sprintf(
		"ims.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(
		avp.DestinationRealm,
		avp.Mbit,
		0,
		datatype.UTF8String(destRealmVal), // or DiameterIdentity
	)

	// Accounting-Record-Type=4 => e.g. EVENT_RECORD or STOP
	// In standard Diameter, values are: 1=Start, 2=Interim, 3=Stop, 4=Event
	artAVP := avp.AccountingRecordType   // 480
	arnAVP := avp.AccountingRecordNumber // 485
	m.NewAVP(artAVP, avp.Mbit, 0, datatype.Enumerated(4))
	m.NewAVP(arnAVP, avp.Mbit, 0, datatype.Unsigned32(0))

	// User-Name => sip:<phoneNumberCalling>...
	userNameVal := fmt.Sprintf(
		"sip:%s@ims.mnc%s.mcc%s.3gppnetwork.org",
		phoneNumberCalling, cfg.MNC, cfg.MCC,
	)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userNameVal))

	// Service-Context-Id => ext.02.001.8.32260@3gpp.org
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("ext.02.001.8.32260@3gpp.org"),
	)

	// 2) Service-Information => grouped
	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// SubscriptionId => type=2 => SIP
			diam.NewAVP(
				avp.SubscriptionID,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(2)),
						diam.NewAVP(
							avp.SubscriptionIDData,
							avp.Mbit,
							0,
							datatype.UTF8String(userNameVal),
						),
					},
				},
			),
			// IMSInformation => grouped
			diam.NewAVP(
				avp.IMSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// EventType => grouped
						diam.NewAVP(
							avp.EventType,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// SIPMethod => "dummy"
									diam.NewAVP(
										avp.SIPMethod,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.UTF8String("dummy"),
									),
									// Event => "dummy"
									diam.NewAVP(
										avp.Event,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.UTF8String("dummy"),
									),
								},
							},
						),
						// RoleOfNode=0 => (Originating)
						diam.NewAVP(
							avp.RoleOfNode,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(0),
						),
						// NodeFunctionality=0
						diam.NewAVP(
							avp.NodeFunctionality,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(0),
						),
						// UserSessionId => "tTOM6k7fswFpIZDJ3qy90g..@10.46.0.3"
						diam.NewAVP(
							avp.UserSessionID,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String("tTOM6k7fswFpIZDJ3qy90g..@10.46.0.3"),
						),
						// CallingPartyAddress => sip:<phoneNumberCalling>...
						diam.NewAVP(
							avp.CallingPartyAddress,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String(userNameVal),
						),
						// CalledPartyAddress => tel:<phoneNumberCalled>
						diam.NewAVP(
							avp.CalledPartyAddress,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf("tel:%s", phoneNumberCalled)),
						),
						// TrunkGroupId => grouped
						diam.NewAVP(
							avp.TrunkGroupID,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(
										avp.OutgoingTrunkGroupID,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.Unsigned32(0),
									),
									diam.NewAVP(
										avp.IncomingTrunkGroupID,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.Unsigned32(0),
									),
								},
							},
						),
						// AccessNetworkInformation => "3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=<cfg.MCC>200001000010b"
						diam.NewAVP(
							avp.AccessNetworkInformation,
							avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf(
								"3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=%s200001000010b",
								cfg.MCC,
							)),
						),
						// TimeStamps => grouped => SIPRequestTimestamp => now
						diam.NewAVP(
							avp.TimeStamps,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(
										avp.SIPRequestTimestamp,
										avp.Vbit,
										Abbas,
										datatype.Time(time.Now()),
									),
								},
							},
						),
					},
				},
			),
		},
	}
	m.NewAVP(avp.ServiceInformation, avp.Mbit|avp.Vbit, Abbas, serviceInfoGrouped)

	// 3) VendorSpecificApplicationId => grouped
	vsaiGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// VendorId=10415
			diam.NewAVP(Abbas, avp.Mbit, 0, datatype.Unsigned32(Abbas)),
			// AuthApplicationId=4
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Enumerated(4)),
		},
	}
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, vsaiGrouped)

	// 4) CC-Request-Type=3 => TERMINATE_REQUEST
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(3))

	// CC-Request-Number=9
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(9))

	// Event-Timestamp => now
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 5) UserEquipmentInfo => grouped
	m.NewAVP(
		avp.UserEquipmentInfo,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				// 0 => IMEI
				diam.NewAVP(
					avp.UserEquipmentInfoType,
					avp.Mbit,
					0,
					datatype.Enumerated(0),
				),
				diam.NewAVP(
					avp.UserEquipmentInfoValue,
					avp.Mbit,
					0,
					datatype.OctetString("doops"),
				),
			},
		},
	)

	// 6) Another SubscriptionId => type=2 => SIP
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(2)),
				diam.NewAVP(
					avp.SubscriptionIDData,
					avp.Mbit,
					0,
					datatype.UTF8String(userNameVal),
				),
			},
		},
	)

	// 7) MultipleServicesIndicator=1 => MULTIPLE_SERVICES_SUPPORTED
	m.NewAVP(avp.MultipleServicesIndicator, avp.Mbit, 0, datatype.Enumerated(1))

	// 8) MultipleServicesCreditControl => grouped
	//    Contains only UsedServiceUnit => { CC-Time=voiceCallUsedTime }, plus ServiceIdentifier & RatingGroup
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// UsedServiceUnit => grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(
							avp.CCTime,
							avp.Mbit,
							0,
							datatype.Unsigned32(5),
						),
					},
				},
			),
			// ServiceIdentifier=1000
			diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(1000)),
			// RatingGroup=100
			diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(100)),
		},
	}
	m.NewAVP(
		avp.MultipleServicesCreditControl,
		avp.Mbit,
		0,
		msccGrouped,
	)

	// 9) TerminationCause=1 => e.g., DIAM_USER_REQUESTED
	m.NewAVP(
		avp.TerminationCause,
		avp.Mbit,
		0,
		datatype.Enumerated(1),
	)

	// Return the built CCR
	return m
}

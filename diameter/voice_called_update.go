package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildVoiceCalledUpdateSessionCCR parallels your Node.js "voiceCalledUpdate" function in Go.
func BuildVoiceCalledUpdateSessionCCR(
	sessionID string,
	phoneNumberCalling string,
	phoneNumberCalled string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// Command Code=272 (Credit-Control), App ID=4 (commonly used for Gx or IMS).
	// Your snippet uses CCR with Accounting-Record-Type, which some operators do in IMS flows.
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 1) Standard AVPs

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
		datatype.DiameterIdentity(originHostVal),
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
		datatype.DiameterIdentity(originRealmVal),
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
		datatype.DiameterIdentity(destRealmVal),
	)

	// Accounting-Record-Type=3 => STOP; Accounting-Record-Number=0
	artAVP := avp.AccountingRecordType   // 480
	arnAVP := avp.AccountingRecordNumber // 485
	m.NewAVP(artAVP, avp.Mbit, 0, datatype.Enumerated(3))
	m.NewAVP(arnAVP, avp.Mbit, 0, datatype.Unsigned32(0))

	// User-Name => phoneNumberCalled
	m.NewAVP(
		avp.UserName,
		avp.Mbit,
		0,
		datatype.UTF8String(phoneNumberCalled),
	)

	// Service-Context-Id
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("ext.02.001.8.32260@3gpp.org"),
	)

	// 2) ServiceInformation => grouped
	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Subscription-Id => Type=0 (END_USER_E164?), Data=phoneNumberCalled
			diam.NewAVP(
				avp.SubscriptionID,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
						diam.NewAVP(
							avp.SubscriptionIDData,
							avp.Mbit,
							0,
							datatype.UTF8String(phoneNumberCalled),
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
						// RoleOfNode=1 (Terminating side)
						diam.NewAVP(
							avp.RoleOfNode,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(1),
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
						// CallingPartyAddress => sip:<phoneNumberCalling>@ims...
						diam.NewAVP(
							avp.CallingPartyAddress,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf(
								"sip:%s@ims.mnc%s.mcc%s.3gppnetwork.org",
								phoneNumberCalling, cfg.MNC, cfg.MCC,
							)),
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
						// Access-Network-Information
						diam.NewAVP(
							avp.AccessNetworkInformation,
							avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf(
								"3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=%s200001000010b",
								cfg.MCC,
							)),
						),
						// TimeStamps => grouped
						diam.NewAVP(
							avp.TimeStamps,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// SIPRequestTimestamp => now
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

	// 4) CC-Request-Type=2 (UPDATE_REQUEST), CC-Request-Number=2
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(2))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(2))

	// 5) Event-Timestamp => now
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 6) UserEquipmentInfo => grouped
	m.NewAVP(
		avp.UserEquipmentInfo,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				// Type=0 => IMEI
				diam.NewAVP(avp.UserEquipmentInfoType, avp.Mbit, 0, datatype.Enumerated(0)),
				diam.NewAVP(
					avp.UserEquipmentInfoValue,
					avp.Mbit,
					0,
					datatype.OctetString("doops"),
				),
			},
		},
	)

	// 7) Another SubscriptionId => Type=0 => phoneNumberCalled
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
				diam.NewAVP(
					avp.SubscriptionIDData,
					avp.Mbit,
					0,
					datatype.UTF8String(phoneNumberCalled),
				),
			},
		},
	)

	// 8) MultipleServicesIndicator=1 => MULTIPLE_SERVICES_SUPPORTED
	m.NewAVP(avp.MultipleServicesIndicator, avp.Mbit, 0, datatype.Enumerated(1))

	// 9) MultipleServicesCreditControl => grouped
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// RequestedServiceUnit => grouped
			diam.NewAVP(
				avp.RequestedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// CCTime => VoiceRequestedTime
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
			diam.NewAVP(
				avp.ServiceIdentifier,
				avp.Mbit,
				0,
				datatype.Unsigned32(1000),
			),
			// RatingGroup=100
			diam.NewAVP(
				avp.RatingGroup,
				avp.Mbit,
				0,
				datatype.Unsigned32(100),
			),
			// UsedServiceUnit => grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// CCTime => VoiceCallUsedTime
						diam.NewAVP(
							avp.CCTime,
							avp.Mbit,
							0,
							datatype.Unsigned32(5),
						),
					},
				},
			),
		},
	}
	m.NewAVP(
		avp.MultipleServicesCreditControl,
		avp.Mbit,
		0,
		msccGrouped,
	)

	// Return the built CCR
	return m
}

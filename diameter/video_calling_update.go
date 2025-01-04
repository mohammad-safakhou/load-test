package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildVideoCallingUpdateSessionCCR parallels your Node.js "videoCallingUpdate" function in Go.
func BuildVideoCallingUpdateSessionCCR(
	sessionID string,
	phoneNumberCalling string,
	phoneNumberCalled string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}
	// Placeholder constants for IMS/3GPP-specific AVPs
	// Adjust them to match your dictionary.

	// Create a CCR message: Command-Code = 272 (Credit-Control), App ID = 4.
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 1) Basic / Standard AVPs

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

	// Accounting-Record-Type=3 => STOP, Accounting-Record-Number=0
	artAVP := avp.AccountingRecordType   // 480
	arnAVP := avp.AccountingRecordNumber // 485
	m.NewAVP(artAVP, avp.Mbit, 0, datatype.Enumerated(3))
	m.NewAVP(arnAVP, avp.Mbit, 0, datatype.Unsigned32(0))

	// User-Name => sip:<phoneNumberCalling>...
	userNameVal := fmt.Sprintf(
		"sip:%s@ims.mnc%s.mcc%s.3gppnetwork.org",
		phoneNumberCalling, cfg.MNC, cfg.MCC,
	)
	m.NewAVP(
		avp.UserName,
		avp.Mbit,
		0,
		datatype.UTF8String(userNameVal),
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
			// SubscriptionId => type=2 => sip
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
						// RoleOfNode=0 (originating side)
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
						// AccessNetworkInformation => "3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=${cfg[0].get.MCC}200001000010b"
						// We'll just inline cfg.MCC below
						diam.NewAVP(
							avp.AccessNetworkInformation,
							avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf(
								"3GPP-E-UTRAN-FDD;utran-cell-id-3gpp=%s200001000010b",
								cfg.MCC,
							)),
						),
						// TimeStamps => grouped (SIPRequestTimestamp => now)
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
			diam.NewAVP(
				Abbas,
				avp.Mbit,
				0,
				datatype.Unsigned32(Abbas),
			),
			diam.NewAVP(
				avp.AuthApplicationID,
				avp.Mbit,
				0,
				datatype.Enumerated(4),
			),
		},
	}
	m.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, vsaiGrouped)

	// 4) CC-Request-Type=2 (UPDATE_REQUEST), CC-Request-Number=1
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(2))
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(1))

	// 5) Event-Timestamp => now
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 6) UserEquipmentInfo => grouped
	m.NewAVP(
		avp.UserEquipmentInfo,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				// 0 => IMEI
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

	// 7) Another SubscriptionId => type=2 => sip
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(
					avp.SubscriptionIDType,
					avp.Mbit,
					0,
					datatype.Enumerated(2), // e.g., sip
				),
				diam.NewAVP(
					avp.SubscriptionIDData,
					avp.Mbit,
					0,
					datatype.UTF8String(userNameVal),
				),
			},
		},
	)

	// 8) MultipleServicesIndicator=1 => MULTIPLE_SERVICES_SUPPORTED
	m.NewAVP(avp.MultipleServicesIndicator, avp.Mbit, 0, datatype.Enumerated(1))

	// 9) MultipleServicesCreditControl => grouped
	//    includes RequestedServiceUnit (CCTime=videoRequestedTime),
	//    ServiceIdentifier=1001, RatingGroup=200,
	//    UsedServiceUnit (CCTime=videoUsedTime)
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// RequestedServiceUnit => grouped
			diam.NewAVP(
				avp.RequestedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(
							avp.CCTime,
							avp.Mbit,
							0,
							datatype.Unsigned32(cfg.VideoRequestedTime),
						),
					},
				},
			),
			// ServiceIdentifier=1001
			diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(1001)),
			// RatingGroup=200
			diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(200)),

			// UsedServiceUnit => grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// CC-Time => videoUsedTime
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

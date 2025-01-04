package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildVideoCallingInitSessionCCR builds a CCR message for a Video Calling scenario.
func BuildVideoCallingInitSessionCCR(
	sessionID string,
	phoneNumberCalling string,
	phoneNumberCalled string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// For VendorSpecificApplicationId grouping (commonly code=260)
	const (
		VendorSpecificApplicationID = avp.VendorSpecificApplicationID
	)

	// Accounting-Record-Type (480) => enumerated(2) for INTERIM
	// Accounting-Record-Number (485) => e.g., 0
	artAVP := avp.AccountingRecordType   // 480
	arnAVP := avp.AccountingRecordNumber // 485

	// ----------------------------------------------------------------
	// Create a new CCR message
	//   Command-Code = 272 (Credit-Control),
	//   Application-ID = 4 (commonly Gx or IMS usage—adjust as needed).
	// ----------------------------------------------------------------
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 1) Basic AVPs: Session-Id, Origin-Host, Origin-Realm, Destination-Realm, etc.

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
	// (No Destination-Host in this snippet, so we skip it.)
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

	// 2) Accounting-Record-Type = 2, Accounting-Record-Number = 0
	m.NewAVP(artAVP, avp.Mbit, 0, datatype.Enumerated(2))
	m.NewAVP(arnAVP, avp.Mbit, 0, datatype.Unsigned32(0))

	// 3) User-Name
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

	// 4) Service-Context-Id
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("ext.02.001.8.32260@3gpp.org"),
	)

	// 5) Service-Information (3GPP vendor, grouped)
	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Subscription-Id
			diam.NewAVP(
				avp.SubscriptionID,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						diam.NewAVP(
							avp.SubscriptionIDType,
							avp.Mbit,
							0,
							datatype.Enumerated(2), // 2 => e.g., END_USER_SIP_URI
						),
						diam.NewAVP(
							avp.SubscriptionIDData,
							avp.Mbit,
							0,
							datatype.UTF8String(userNameVal),
						),
					},
				},
			),
			// IMS-Information (vendor-specific)
			diam.NewAVP(
				avp.IMSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// Event-Type (grouped)
						diam.NewAVP(
							avp.EventType,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// SIPMethod="INVITE"
									diam.NewAVP(
										avp.SIPMethod,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.UTF8String("INVITE"),
									),
									// Expires=4294967295
									diam.NewAVP(
										avp.Expires,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.Unsigned32(4294967295),
									),
								},
							},
						),
						// Role-Of-Node=0
						diam.NewAVP(
							avp.RoleOfNode,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(0),
						),
						// Node-Functionality=0
						diam.NewAVP(
							avp.NodeFunctionality,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(0),
						),
						// User-Session-Id
						diam.NewAVP(
							avp.UserSessionID,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String("tTOM6k7fswFpIZDJ3qy90g..@10.46.0.3"),
						),
						// Calling-Party-Address
						diam.NewAVP(
							avp.CallingPartyAddress,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String(userNameVal),
						),
						// Called-Party-Address
						diam.NewAVP(
							avp.CalledPartyAddress,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.UTF8String(fmt.Sprintf("tel:%s", phoneNumberCalled)),
						),
						// Trunk-Group-Id
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
						// Time-Stamps (grouped)
						diam.NewAVP(
							avp.TimeStamps,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// SIP-Request-Timestamp
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

	m.NewAVP(
		avp.ServiceInformation,
		avp.Mbit|avp.Vbit,
		Abbas,
		serviceInfoGrouped,
	)

	// 6) Vendor-Specific-Application-Id
	vsaiGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Vendor-Id=10415
			diam.NewAVP(Abbas, avp.Mbit, 0, datatype.Unsigned32(10415)),
			// Auth-Application-Id=4
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Enumerated(4)),
		},
	}
	m.NewAVP(
		VendorSpecificApplicationID,
		avp.Mbit,
		0,
		vsaiGrouped,
	)

	// 7) CC-Request-Type=1 (INITIAL_REQUEST)
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1))

	// 8) CC-Request-Number=0
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))

	// 9) Event-Timestamp
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 10) User-Equipment-Info (AVP=458 in base, but confirm with your dictionary)
	m.NewAVP(
		avp.UserEquipmentInfo,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(
					avp.UserEquipmentInfoType,
					avp.Mbit,
					0,
					datatype.Enumerated(0), // 0 => IMEI
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

	// 11) Additional Subscription-Id
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

	// 12) Multiple-Services-Indicator = 1 (MULTIPLE_SERVICES_SUPPORTED)
	m.NewAVP(
		avp.MultipleServicesIndicator,
		avp.Mbit,
		0,
		datatype.Enumerated(1),
	)

	// 13) Multiple-Services-Credit-Control (456)
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Requested-Service-Unit (437) => includes CC-Time (420)
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
							datatype.Unsigned32(5),
						),
					},
				},
			),
			// Service-Identifier (439)
			diam.NewAVP(
				avp.ServiceIdentifier,
				avp.Mbit,
				0,
				datatype.Unsigned32(1001),
			),
			// Rating-Group (432)
			diam.NewAVP(
				avp.RatingGroup,
				avp.Mbit,
				0,
				datatype.Unsigned32(200),
			),
		},
	}
	m.NewAVP(
		avp.MultipleServicesCreditControl,
		avp.Mbit,
		0,
		msccGrouped,
	)

	// Return the fully built CCR
	return m
}

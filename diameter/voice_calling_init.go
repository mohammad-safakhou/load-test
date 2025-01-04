package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildVoiceCallingInitSessionCCR builds a CCR request for IMS voice calling.
func BuildVoiceCallingInitSessionCCR(
	sessionID string,
	phoneNumberCalling string,
	phoneNumberCalled string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// -----------------------------

	// For VendorSpecificApplicationId grouping
	// (commonly used with IMS, might be AVP 260 in base RFC)
	// Also, note for IMS your Auth-Application-Id might be 5, 6, or 7 for different interfaces,
	// or you might keep it at 4 if it’s aligned with a 3GPP app ID.
	const (
		VendorSpecificApplicationID = avp.VendorSpecificApplicationID // 260
	)

	// Accounting-Record-Type codes (AVP=480) for IMS: 1=Start, 2=Interim, 3=Stop, etc.
	// This example sets it to enumerated(2).
	// If you’re using CC-Request-Type, that’s also okay. Some IMS flows use Accounting-Record-Type
	// plus CC-Request-Type in a combined message.
	// Adjust as needed for your flow.

	// --------------------------------------------------------------------
	// Create a new CCR message
	//   - Typically, you might see diameter.Accounting, or diameter.CreditControl
	//   - For a typical IMS “Rf” interface, you might do diam.Accounting with AppID=3 or 4, etc.
	//   - For an IMS-based call detail flow, the command code is often 272 (Credit-Control)
	//     or 4xx if using accounting. Adjust to your scenario.
	// --------------------------------------------------------------------
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// 1) Add basic AVPs (Session-Id, Origin-Host, Realm, etc.)
	sessionIDVal := fmt.Sprintf("smf.epc.mnc%s.mcc%s.3gppnetwork.org;%s;%s",
		cfg.MNC, cfg.MCC, sessionID, sessionID)
	m.NewAVP(
		avp.SessionID,
		avp.Mbit,
		0,
		datatype.UTF8String(sessionIDVal),
	)

	// Origin-Host
	originHostVal := fmt.Sprintf("scscf.ims.mnc%s.mcc%s.3gppnetwork.org", cfg.MNC, cfg.MCC)
	m.NewAVP(
		avp.OriginHost,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(originHostVal),
	)

	// Origin-Realm
	originRealmVal := fmt.Sprintf("ims.mnc%s.mcc%s.3gppnetwork.org", cfg.MNC, cfg.MCC)
	m.NewAVP(
		avp.OriginRealm,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(originRealmVal),
	)

	// Destination-Host
	// Example: "hssocs.voiceblue.com"
	m.NewAVP(
		avp.DestinationHost,
		avp.Mbit,
		0,
		datatype.DiameterIdentity("hssocs.voiceblue.com"),
	)

	// Destination-Realm
	destRealmVal := fmt.Sprintf("ims.mnc%s.mcc%s.3gppnetwork.org", cfg.MNC, cfg.MCC)
	m.NewAVP(
		avp.DestinationRealm,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(destRealmVal),
	)

	// 2) Add Accounting or CCR details

	// Accounting-Record-Type (AVP=480 typically) -> enumerated(2) for INTERIM
	// If you have the code as part of standard dict, you can do:
	//   diam.NewAVP(avp.AccountingRecordType, avp.Mbit, 0, datatype.Enumerated(2))
	// Or define your own if it’s not in your dictionary
	artAVP := avp.AccountingRecordType // 480 in base
	m.NewAVP(artAVP, avp.Mbit, 0, datatype.Enumerated(2))

	// Accounting-Record-Number (AVP=485)
	arnAVP := avp.AccountingRecordNumber // 485 in base
	m.NewAVP(arnAVP, avp.Mbit, 0, datatype.Unsigned32(0))

	// User-Name (AVP=1)
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
	// For IMS, you had "ext.02.001.8.32260@3gpp.org"
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("ext.02.001.8.32260@3gpp.org"),
	)

	// 3) Service-Information (3GPP vendor, grouped)
	// Code commonly 873 or custom. We pass V-bit + M-bit with vendor=10415.
	// Inside it: Subscription-Id + IMS-Information
	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Subscription-Id (AVP=443 in base; or custom if vendor-specific)
			diam.NewAVP(
				avp.SubscriptionID,
				avp.Mbit,
				0, // no vendor ID if standard
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// Subscription-Id-Type (AVP=450), enumerated(2) => END_USER_SIP_URI?
						diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(2)),
						// Subscription-Id-Data (AVP=444)
						diam.NewAVP(
							avp.SubscriptionIDData,
							avp.Mbit,
							0,
							datatype.UTF8String(userNameVal),
						),
					},
				},
			),

			// IMS-Information (vendor-specific code, e.g., 876 in our placeholders)
			diam.NewAVP(
				avp.IMSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// Event-Type
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
									// OutgoingTrunkGroupId=0
									diam.NewAVP(
										avp.OutgoingTrunkGroupID,
										avp.Mbit|avp.Vbit,
										Abbas,
										datatype.Unsigned32(0),
									),
									// IncomingTrunkGroupId=0
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
						// TimeStamps
						diam.NewAVP(
							avp.TimeStamps,
							avp.Mbit|avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// SIPRequestTimestamp
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

	// Add Service-Information with V & M bits, vendor=10415
	m.NewAVP(
		avp.ServiceInformation,
		avp.Mbit|avp.Vbit,
		Abbas,
		serviceInfoGrouped,
	)

	// 4) VendorSpecificApplicationId (AVP=260)
	vsaiGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Vendor-Id (AVP=266)
			diam.NewAVP(Abbas, avp.Mbit, 0, datatype.Unsigned32(10415)),
			// Auth-Application-Id => 4 (IMS might use 4, 5, or others)
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Enumerated(4)),
		},
	}
	m.NewAVP(
		VendorSpecificApplicationID,
		avp.Mbit,
		0,
		vsaiGrouped,
	)

	// 5) CC-Request-Type (416)=1 => INITIAL_REQUEST
	m.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Enumerated(1))

	// CC-Request-Number (415)=0
	m.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))

	// Event-Timestamp (55)
	m.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))

	// 6) User-Equipment-Info (AVP=458, standard, but ensure your dictionary is correct)
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
					datatype.Enumerated(0), // 0=IMEI
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

	// 7) Another Subscription-Id grouped (for SIP)
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

	// 8) MultipleServicesIndicator (AVP=455)
	m.NewAVP(
		avp.MultipleServicesIndicator,
		avp.Mbit,
		0,
		datatype.Enumerated(1), // 1=MULTIPLE_SERVICES_SUPPORTED
	)

	// 9) MultipleServicesCreditControl (AVP=456)
	//    This includes RequestedServiceUnit, ServiceIdentifier, RatingGroup, etc.
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// RequestedServiceUnit (AVP=437)
			diam.NewAVP(
				avp.RequestedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// CC-Time (AVP=420)
						diam.NewAVP(
							avp.CCTime,
							avp.Mbit,
							0,
							datatype.Unsigned32(5),
						),
					},
				},
			),
			// Service-Identifier (AVP=439)
			diam.NewAVP(
				avp.ServiceIdentifier,
				avp.Mbit,
				0,
				datatype.Unsigned32(1000),
			),
			// Rating-Group (AVP=432)
			diam.NewAVP(
				avp.RatingGroup,
				avp.Mbit,
				0,
				datatype.Unsigned32(100),
			),
		},
	}
	m.NewAVP(
		avp.MultipleServicesCreditControl,
		avp.Mbit,
		0,
		msccGrouped,
	)

	// Done! Return the fully built CCR
	return m
}

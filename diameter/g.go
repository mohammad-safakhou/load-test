package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// Example structures to hold config/user data:
type Config struct {
	MNC       string
	MCC       string
	BitrateUL uint32
	BitrateDL uint32
}

type UserData struct {
	PreMiddleFix     string
	Prefix           string
	TimeZone         string
	UserLocationInfo string
}

// Example function that returns a CCR (Credit-Control Request) message
func BuildInitSessionCCR(
	sessionID string,
	phoneNumber string,
	cfg *Config,
	userData *UserData,
	imeiMap map[string]string, // e.g., IMEISVs
) *diam.Message {

	// --------------------------------------------------------------------------------
	// 1) Create a new CCR message.
	//    - 272 is the command code for Credit-Control.
	//    - 4 is often the Application ID for 3GPP Gx/Gy; change if you have a custom app.
	// --------------------------------------------------------------------------------
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	// --------------------------------------------------------------------------------
	// 2) Add AVPs to the CCR message.
	//    Each line here corresponds to your original JavaScript example.
	// --------------------------------------------------------------------------------

	// Session-Id
	// e.g., "smf.epc.mnc999.mcc999.3gppnetwork.org;sessionID;sessionID"
	m.NewAVP(
		avp.SessionID,
		avp.Mbit,
		0,
		datatype.UTF8String(fmt.Sprintf(
			"smf.epc.mnc%s.mcc%s.3gppnetwork.org;%s;%s",
			cfg.MNC, cfg.MCC, sessionID, sessionID,
		)),
	)

	// Origin-Host
	m.NewAVP(
		avp.OriginHost,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(fmt.Sprintf(
			"smf.epc.mnc%s.mcc%s.3gppnetwork.org",
			cfg.MNC, cfg.MCC,
		)),
	)

	// Origin-Realm
	m.NewAVP(
		avp.OriginRealm,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(fmt.Sprintf(
			"epc.mnc%s.mcc%s.3gppnetwork.org",
			cfg.MNC, cfg.MCC,
		)),
	)

	// Destination-Realm
	m.NewAVP(
		avp.DestinationRealm,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(fmt.Sprintf(
			"epc.mnc%s.mcc%s.3gppnetwork.org",
			cfg.MNC, cfg.MCC,
		)),
	)

	// Auth-Application-Id (4 for Gx)
	m.NewAVP(
		avp.AuthApplicationID,
		avp.Mbit,
		0,
		datatype.Unsigned32(4),
	)

	// Service-Context-Id (3GPP uses "32251@3gpp.org" for Gx)
	// 461 is the standard AVP code for Service-Context-Id.
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("32251@3gpp.org"),
	)

	// CC-Request-Type (AVP code 416; 1=INITIAL_REQUEST)
	m.NewAVP(
		avp.CCRequestType,
		avp.Mbit,
		0,
		datatype.Enumerated(1), // 1 = INITIAL_REQUEST
	)

	// CC-Request-Number (AVP code 415)
	m.NewAVP(
		avp.CCRequestNumber,
		avp.Mbit,
		0,
		datatype.Unsigned32(0),
	)

	// Event-Timestamp (AVP code 55)
	m.NewAVP(
		avp.EventTimestamp,
		avp.Mbit,
		0,
		datatype.Time(time.Now()),
	)

	// Subscription-Id (AVP code 443) #1
	// Type=0 (IMS subscription, e.g. END_USER_E164)
	// Data=phoneNumber
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(phoneNumber)),
		}},
	)

	// Subscription-Id (AVP code 443) #2
	// Type=1 (END_USER_IMSI, for instance)
	// Data=userData.PreMiddleFix + phoneNumber
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(1)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(
				userData.PreMiddleFix+phoneNumber,
			)),
		}},
	)

	// Subscription-Id (AVP code 443) #3
	// Repeated as in your example (though typically you wouldn't repeat the same type=0)
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(phoneNumber)),
		}},
	)

	// Requested-Action (AVP code 436; 0 = CHECK_BALANCE ?)
	m.NewAVP(
		avp.RequestedAction,
		avp.Mbit,
		0,
		datatype.Enumerated(0),
	)

	// AoC-Request-Type (Vendor-Specific)
	// If it’s 3GPP vendor ID (10415), you need avp.Vbit along with Mbit
	// For demonstration, using AVP code 812 for AoC-Request-Type (may differ in your dictionary)
	const TGPPVendorID = 10415
	const AoCRequestTypeAVP = 812 // Example code; adjust for your dictionary
	m.NewAVP(
		AoCRequestTypeAVP,
		avp.Mbit|avp.Vbit,
		TGPPVendorID,
		datatype.Enumerated(1), // E.g., AOC request type
	)

	// Multiple-Services-Indicator (AVP code 455; 0 = MULTIPLE_SERVICES_NOT_SUPPORTED)
	m.NewAVP(
		avp.MultipleServicesIndicator,
		avp.Mbit,
		0,
		datatype.Enumerated(0),
	)

	// Multiple-Services-Credit-Control (AVP code 456)
	// Includes Requested-Service-Unit (code 437), Used-Service-Unit (446), QoS-Information, etc.
	// This is a grouped AVP containing sub-AVPs
	// Many of these codes (QoSInformation, TGPPRATType, etc.) are 3GPP-specific and may
	// require custom dictionary definitions.
	const QoSInformationAVP = 1016
	const TGPPRATTypeAVP = 1032

	msccGroupedAVP := &diam.GroupedAVP{AVP: []*diam.AVP{
		// Requested-Service-Unit (437) -> often contains CC-Time, CC-Total-Octets, etc. Here it’s empty
		diam.NewAVP(avp.RequestedServiceUnit, avp.Mbit, 0, datatype.UTF8String("")),

		// Used-Service-Unit (446) -> sub-group for usage
		diam.NewAVP(avp.UsedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{AVP: []*diam.AVP{
			diam.NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(0)),
			diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(0)),
			diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(0)),
		}}),

		// QoS-Information (3GPP-specific)
		diam.NewAVP(QoSInformationAVP, avp.Mbit|avp.Vbit, TGPPVendorID, &diam.GroupedAVP{AVP: []*diam.AVP{
			// QoS-Class-Identifier = 1028
			diam.NewAVP(1028, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Enumerated(5)), // e.g., QCI = 5
			// Allocation-Retention-Priority = 1034
			diam.NewAVP(1034, avp.Vbit, TGPPVendorID, &diam.GroupedAVP{AVP: []*diam.AVP{
				// Priority-Level = 1046
				diam.NewAVP(1046, avp.Vbit, TGPPVendorID, datatype.Enumerated(1)),
				// Preemption-Capability = 1047
				diam.NewAVP(1047, avp.Vbit, TGPPVendorID, datatype.Enumerated(1)),
				// Preemption-Vulnerability = 1048
				diam.NewAVP(1048, avp.Vbit, TGPPVendorID, datatype.Enumerated(1)),
			}}),
			// APN-Aggregate-Max-Bitrate-UL = 1041
			diam.NewAVP(1041, avp.Vbit, TGPPVendorID, datatype.Unsigned32(cfg.BitrateUL)),
			// APN-Aggregate-Max-Bitrate-DL = 1042
			diam.NewAVP(1042, avp.Vbit, TGPPVendorID, datatype.Unsigned32(cfg.BitrateDL)),
		}}),

		// TGPP-RAT-Type (code 1032?), pass "06" as OctetString
		diam.NewAVP(TGPPRATTypeAVP, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.OctetString("06")),
	}}

	m.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, msccGroupedAVP)

	// Service-Information (3GPP vendor-specific)
	// 873 or something similar in your dictionary
	const ServiceInformationAVP = 873
	const PSInformationAVP = 874

	psInfoGroupedAVP := &diam.GroupedAVP{AVP: []*diam.AVP{
		// TGPP-Charging-ID, 2? (One example; actual code depends on dictionary)
		diam.NewAVP(2, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.OctetString("0000000a")),

		// TGPP-PDP-Type
		// Some references define it around 3 or 5. You’ll need a custom dictionary for real code.
		diam.NewAVP(3, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Enumerated(0)),

		// PDP-Address
		// Might be AVP code 1227 for Gx
		diam.NewAVP(1227, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Address("10.46.0.2")),

		// SGSN-Address
		diam.NewAVP(1228, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Address("127.0.0.3")),

		// GGSN-Address (repeated)
		diam.NewAVP(1229, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Address("127.0.0.4")),
		diam.NewAVP(1229, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.Address("127.0.0.4")),

		// Called-Station-Id (30)
		diam.NewAVP(avp.CalledStationID, avp.Mbit, 0, datatype.UTF8String("internet")),

		// TGPP-Selection-Mode?
		diam.NewAVP(4, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.UTF8String("0")),

		// TGPP-SGSN-MCCMNC
		// Note: Typically code for TGPP-SGSN-MCCMNC might differ.
		diam.NewAVP(18, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.UTF8String(userData.Prefix)),

		// TGPP-NSAPI
		diam.NewAVP(5, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.OctetString("\\005")),

		// TGPP-MS-TimeZone
		diam.NewAVP(23, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.OctetString(userData.TimeZone)),

		// TGPP-User-Location-Info
		diam.NewAVP(22, avp.Mbit|avp.Vbit, TGPPVendorID, datatype.UTF8String(userData.UserLocationInfo)),

		// User-Equipment-Info (code 458 in RFC, but 3GPP also extends it)
		diam.NewAVP(avp.UserEquipmentInfo, avp.Mbit, 0, &diam.GroupedAVP{AVP: []*diam.AVP{
			diam.NewAVP(avp.UserEquipmentInfoType, avp.Mbit, 0, datatype.Enumerated(0)), // 0=IMEI
			diam.NewAVP(avp.UserEquipmentInfoValue, avp.Mbit, 0, datatype.OctetString(imeiMap[sessionID])),
		}}),
	}}

	serviceInfoGroupedAVP := &diam.GroupedAVP{AVP: []*diam.AVP{
		diam.NewAVP(PSInformationAVP, avp.Mbit|avp.Vbit, TGPPVendorID, psInfoGroupedAVP),
	}}

	m.NewAVP(
		ServiceInformationAVP,
		avp.Mbit|avp.Vbit,
		TGPPVendorID,
		serviceInfoGroupedAVP,
	)

	// --------------------------------------------------------------------------------
	// 3) Return the fully built CCR message.
	// --------------------------------------------------------------------------------
	return m
}

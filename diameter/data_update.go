package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

type IMEIMap map[string]string

// BuildDataUpdateSessionCCR creates a CCR for data usage update,
// closely mirroring your Node.js snippet.
func BuildDataUpdateSessionCCR(
	sessionID string,
	phoneNumber string,
) *diam.Message {
	cfg := &Config{
		MNC:             "020",
		MCC:             "418",
		UploadedBytes:   53362,
		DownloadedBytes: 5190705,
	}

	// 3GPP Vendor ID

	// Some placeholders for 3GPP-specific AVPs (QoS, RAT type, etc.).
	// These AVP codes often require a custom dictionary if not in dict.Default.

	// Create the CCR message: Command-Code=272 (Credit-Control), App-ID=4
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
		"smf.epc.mnc%s.mcc%s.3gppnetwork.org",
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
		"epc.mnc%s.mcc%s.3gppnetwork.org",
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
		"epc.mnc%s.mcc%s.3gppnetwork.org",
		cfg.MNC, cfg.MCC,
	)
	m.NewAVP(
		avp.DestinationRealm,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(destRealmVal),
	)

	// Auth-Application-Id=4
	m.NewAVP(
		avp.AuthApplicationID,
		avp.Mbit,
		0,
		datatype.Unsigned32(4),
	)

	// Service-Context-Id
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("32251@3gpp.org"),
	)

	// CC-Request-Type=2 (UPDATE_REQUEST)
	m.NewAVP(
		avp.CCRequestType,
		avp.Mbit,
		0,
		datatype.Enumerated(2),
	)

	// CC-Request-Number=1
	m.NewAVP(
		avp.CCRequestNumber,
		avp.Mbit,
		0,
		datatype.Unsigned32(1),
	)

	// Destination-Host
	m.NewAVP(
		avp.DestinationHost,
		avp.Mbit,
		0,
		datatype.DiameterIdentity(fmt.Sprintf(
			"CGR-DA.epc.mnc%s.mcc%s.3gppnetwork.org",
			cfg.MNC, cfg.MCC,
		)),
	)

	// Event-Timestamp => now
	m.NewAVP(
		avp.EventTimestamp,
		avp.Mbit,
		0,
		datatype.Time(time.Now()),
	)

	// 2) Subscription-Id(s)
	// #1 => Type=1 (IMS? or END_USER_IMSI?), Data=userData.PreMiddleFix + phoneNumber
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
					datatype.Enumerated(1),
				),
				diam.NewAVP(
					avp.SubscriptionIDData,
					avp.Mbit,
					0,
					datatype.UTF8String(cfg.PreMiddleFix+phoneNumber),
				),
			},
		},
	)

	// #2 => Type=0 (END_USER_E164?), Data=phoneNumber
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
					datatype.Enumerated(0),
				),
				diam.NewAVP(
					avp.SubscriptionIDData,
					avp.Mbit,
					0,
					datatype.UTF8String(phoneNumber),
				),
			},
		},
	)

	// Requested-Action=0 (CHECK_BALANCE or DIRECT_DEBIT, depends on 3GPP spec)
	m.NewAVP(
		avp.RequestedAction,
		avp.Mbit,
		0,
		datatype.Enumerated(0),
	)

	// AoCRequestType (vendor=10415, with M+V bits)
	m.NewAVP(
		avp.AoCRequestType,
		avp.Mbit|avp.Vbit,
		Abbas,
		datatype.Enumerated(1), // e.g., 1 => AoC request
	)

	// 3) Multiple-Services-Credit-Control => grouped
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// Requested-Service-Unit => empty string in your snippet
			diam.NewAVP(
				avp.RequestedServiceUnit,
				avp.Mbit,
				0,
				datatype.UTF8String(""),
			),
			// Used-Service-Unit => grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// ReportingReason=3 (e.g., AVP=9992 in custom dict)
						diam.NewAVP(
							avp.ReportingReason,
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(3),
						),
						// CC-Time=0
						diam.NewAVP(
							avp.CCTime,
							avp.Mbit,
							0,
							datatype.Unsigned32(0),
						),
						// CC-Input-Octets=downloadedBytes
						diam.NewAVP(
							avp.CCInputOctets,
							avp.Mbit,
							0,
							datatype.Unsigned64(cfg.UploadedBytes),
						),
						// CC-Output-Octets=uploadedBytes
						diam.NewAVP(
							avp.CCOutputOctets,
							avp.Mbit,
							0,
							datatype.Unsigned64(cfg.DownloadedBytes),
						),
					},
				},
			),
			// QoSInformation => grouped
			diam.NewAVP(
				avp.QoSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// QoS-Class-Identifier=9
						diam.NewAVP(
							1028, // e.g., QCI code from 3GPP
							avp.Mbit|avp.Vbit,
							Abbas,
							datatype.Enumerated(9),
						),
						// Allocation-Retention-Priority => grouped
						diam.NewAVP(
							1034, // e.g., ARP code
							avp.Vbit,
							Abbas,
							&diam.GroupedAVP{
								AVP: []*diam.AVP{
									// Priority-Level=8
									diam.NewAVP(1046, avp.Vbit, Abbas, datatype.Enumerated(8)),
									// Preemption-Capability=1
									diam.NewAVP(1047, avp.Vbit, Abbas, datatype.Enumerated(1)),
									// Preemption-Vulnerability=1
									diam.NewAVP(1048, avp.Vbit, Abbas, datatype.Enumerated(1)),
								},
							},
						),
						// APN-Aggregate-Max-Bitrate-UL
						diam.NewAVP(
							1041,
							avp.Vbit,
							Abbas,
							datatype.Unsigned32(cfg.BitrateUL),
						),
						// APN-Aggregate-Max-Bitrate-DL
						diam.NewAVP(
							1040,
							avp.Vbit,
							Abbas,
							datatype.Unsigned32(cfg.BitrateDL),
						),
					},
				},
			),
			// TGPP-RAT-Type => "06" (OctetString)
			diam.NewAVP(
				avp.TGPPRATType,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.OctetString("06"),
			),
		},
	}
	m.NewAVP(
		avp.MultipleServicesCreditControl,
		avp.Mbit,
		0,
		msccGrouped,
	)

	// 4) ServiceInformation => grouped
	psInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// TGPP-Charging-Id, PDP-Address, SGSN/GGSN-Address, etc.
			diam.NewAVP(
				avp.TGPPChargingID,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.OctetString("00000014"),
			),
			// PDP-Address => 10.46.0.8
			diam.NewAVP(
				avp.PDPAddress,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.Address("10.46.0.8"),
			),
			// SGSN-Address => 127.0.0.3
			diam.NewAVP(
				avp.SGSNAddress,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.Address("127.0.0.3"),
			),
			// GGSN-Address => 127.0.0.4 (twice)
			diam.NewAVP(
				avp.GGSNAddress,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.Address("127.0.0.4"),
			),
			diam.NewAVP(
				avp.GGSNAddress,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.Address("127.0.0.4"),
			),
			// Called-Station-Id => "internet"
			diam.NewAVP(
				avp.CalledStationID,
				avp.Mbit,
				0,
				datatype.UTF8String("internet"),
			),
			// TGPP-Selection-Mode => "0"
			diam.NewAVP(
				avp.TGPPSelectionMode,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String("0"),
			),
			// TGPP-SGSN-MCCMNC => userData.Prefix
			diam.NewAVP(
				avp.TGPPSGSNMCCMNC,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String(cfg.Prefix),
			),
			// TGPP-NSAPI => "\006"
			diam.NewAVP(
				avp.TGPPNSAPI,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String("\\006"),
			),
			// TGPP-MS-TimeZone => userData.TimeZone
			diam.NewAVP(
				avp.TGPPMSTimeZone,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.OctetString(cfg.TimeZone),
			),
			// TGPP-User-Location-Info => userData.UserLocationInfo
			diam.NewAVP(
				avp.TGPPUserLocationInfo,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String(cfg.UserLocationInfo),
			),
			// User-Equipment-Info => grouped
			diam.NewAVP(
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
			),
		},
	}

	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// PSInformation => grouped
			diam.NewAVP(
				avp.PSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				psInfoGrouped,
			),
		},
	}
	m.NewAVP(
		avp.ServiceInformation,
		avp.Mbit|avp.Vbit,
		Abbas,
		serviceInfoGrouped,
	)

	// Return the fully built CCR
	return m
}

package diameter

import (
	"fmt"
	"time"

	"github.com/MHG14/go-diameter/v4/diam"
	"github.com/MHG14/go-diameter/v4/diam/avp"
	"github.com/MHG14/go-diameter/v4/diam/datatype"
	"github.com/MHG14/go-diameter/v4/diam/dict"
)

// BuildDataTerminateSessionCCR parallels your Node.js `dataTerminate` function in Go.
func BuildDataTerminateSessionCCR(
	sessionID string,
	phoneNumber string,
) *diam.Message {
	cfg := &Config{
		MNC: "020",
		MCC: "418",
	}

	// 1) Create the CCR message: Command-Code=272 (Credit-Control), App-ID=4
	m := diam.NewRequest(diam.CreditControl, 4, dict.Default)

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

	// Auth-Application-Id=4
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))

	// Service-Context-Id => "32251@3gpp.org"
	m.NewAVP(
		avp.ServiceContextID,
		avp.Mbit,
		0,
		datatype.UTF8String("32251@3gpp.org"),
	)

	// CC-Request-Type=3 => TERMINATE_REQUEST
	m.NewAVP(
		avp.CCRequestType,
		avp.Mbit,
		0,
		datatype.Enumerated(3),
	)

	// CC-Request-Number=1
	m.NewAVP(
		avp.CCRequestNumber,
		avp.Mbit,
		0,
		datatype.Unsigned32(1),
	)

	// Event-Timestamp => now
	m.NewAVP(
		avp.EventTimestamp,
		avp.Mbit,
		0,
		datatype.Time(time.Now()),
	)

	// Subscription-Id #1 => Type=1 (e.g., IMSI?), Data=userData.PreMiddleFix + phoneNumber
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

	// Subscription-Id #2 => Type=0 (END_USER_E164?), Data=phoneNumber
	m.NewAVP(
		avp.SubscriptionID,
		avp.Mbit,
		0,
		&diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(0)),
				diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(phoneNumber)),
			},
		},
	)

	// Termination-Cause=1 (DIAM_USER_REQUESTED?)
	m.NewAVP(avp.TerminationCause, avp.Mbit, 0, datatype.Enumerated(1))

	// Requested-Action=0 => CHECK_BALANCE or DIRECT_DEBIT, depends on spec
	m.NewAVP(
		avp.RequestedAction,
		avp.Mbit,
		0,
		datatype.Enumerated(0),
	)

	// AoCRequestType => enumerated(1), vendor=10415, M+V bits
	m.NewAVP(
		avp.AoCRequestType,
		avp.Mbit|avp.Vbit,
		Abbas,
		datatype.Enumerated(1),
	)

	// MultipleServicesCreditControl => grouped
	msccGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// UsedServiceUnit => grouped
			diam.NewAVP(
				avp.UsedServiceUnit,
				avp.Mbit,
				0,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// CC-Time=0
						diam.NewAVP(avp.CCTime, avp.Mbit, 0, datatype.Unsigned32(0)),
						// CC-Input-Octets => DownloadedBytes
						diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(cfg.DownloadedBytes)),
						// CC-Output-Octets => UploadedBytes
						diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(cfg.UploadedBytes)),
					},
				},
			),
			// Reporting-Reason=2 => e.g., FINAL (?), code from your snippet
			diam.NewAVP(
				avp.ReportingReason,
				avp.Mbit,
				0,
				datatype.Enumerated(2),
			),
			// QoSInformation => grouped
			diam.NewAVP(
				avp.QoSInformation,
				avp.Mbit|avp.Vbit,
				Abbas,
				&diam.GroupedAVP{
					AVP: []*diam.AVP{
						// QoS-Class-Identifier=9
						diam.NewAVP(1028, avp.Mbit|avp.Vbit, Abbas, datatype.Enumerated(9)),
						// Allocation-Retention-Priority => grouped
						diam.NewAVP(1034, avp.Vbit, Abbas, &diam.GroupedAVP{
							AVP: []*diam.AVP{
								// Priority-Level=8
								diam.NewAVP(1046, avp.Vbit, Abbas, datatype.Enumerated(8)),
								// Preemption-Capability=1
								diam.NewAVP(1047, avp.Vbit, Abbas, datatype.Enumerated(1)),
								// Preemption-Vulnerability=1
								diam.NewAVP(1048, avp.Vbit, Abbas, datatype.Enumerated(1)),
							},
						}),
						// APN-Aggregate-Max-Bitrate-UL => cfg.BitrateUL
						diam.NewAVP(1041, avp.Vbit, Abbas, datatype.Unsigned32(cfg.BitrateUL)),
						// APN-Aggregate-Max-Bitrate-DL => cfg.BitrateDL
						diam.NewAVP(1040, avp.Vbit, Abbas, datatype.Unsigned32(cfg.BitrateDL)),
					},
				},
			),
			// TGPPRATType => "06" as OctetString
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

	// ServiceInformation => grouped
	psInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			// TGPP-Charging-Id => "0000000b"
			diam.NewAVP(
				avp.TGPPChargingID,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.OctetString("0000000b"),
			),
			// PDP-Address => 10.45.0.3
			diam.NewAVP(
				avp.PDPAddress,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.Address("10.45.0.3"),
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
			// TGPP-SGSN-MCCMNC => `${cfg.MCC}20`
			diam.NewAVP(
				avp.TGPPSGSNMCCMNC,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String(fmt.Sprintf("%s20", cfg.MCC)),
			),
			// TGPP-NSAPI => "\005"
			diam.NewAVP(
				avp.TGPPNSAPI,
				avp.Mbit|avp.Vbit,
				Abbas,
				datatype.UTF8String("\\005"),
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
			),
		},
	}

	serviceInfoGrouped := &diam.GroupedAVP{
		AVP: []*diam.AVP{
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

	// Return the built CCR
	return m
}

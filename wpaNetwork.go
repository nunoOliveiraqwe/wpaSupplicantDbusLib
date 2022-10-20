package wpaSuppDBusLib

import (
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

type ScanSSID int8
type Mode int8
type Proto string
type AuthAlg string
type KeyManagement string
type PairWise string
type Group string
type EapolFlag int8

const (
	ScanOn                ScanSSID      = 0
	ScanOff               ScanSSID      = 1
	ModeInfrastructure    Mode          = 0
	ModeIBSS              Mode          = 1
	WPAProto              Proto         = "WPA"
	WPA2Proto             Proto         = "RSN"
	WpaEAP                KeyManagement = "WPA-EAP"
	WpaPSK                KeyManagement = "WPA-PSK"
	IEEE8021X             KeyManagement = "IEEE8021X"
	NONE                  KeyManagement = "NONE"
	AuthAlgOpen           AuthAlg       = "OPEN"
	AuthAlgShared         AuthAlg       = "SHARED"
	AuthAlgLeap           AuthAlg       = "LEAP"
	PairWiseCCMP          PairWise      = "CCMP"
	PairWiseTKIP          PairWise      = "TKIP"
	PairWiseNone          PairWise      = "NONE"
	GroupCCMP             Group         = "CCMP"
	GroupTKIP             Group         = "TKIP"
	GroupWEP104           Group         = "WEP104"
	GroupWEP40            Group         = "WEP40"
	EapolDynamicUnicast   EapolFlag     = 1
	EapolDynamicBroadcast EapolFlag     = 2
	EapolDynamicBoth      EapolFlag     = 3
)

var scanSlice = []ScanSSID{ScanOn, ScanOff}
var modeSlice = []Mode{ModeInfrastructure, ModeIBSS}
var protoSlice = []Proto{WPAProto, WPA2Proto}
var keyMngtSlice = []KeyManagement{WpaEAP, WpaPSK, IEEE8021X, NONE}
var authAlgSlice = []AuthAlg{AuthAlgOpen, AuthAlgShared, AuthAlgLeap}
var pairWiseSlice = []PairWise{PairWiseCCMP, PairWiseTKIP, PairWiseNone}
var groupSlice = []Group{GroupCCMP, GroupTKIP, GroupWEP104, GroupWEP40}
var eapFlagSlice = []EapolFlag{EapolDynamicUnicast, EapolDynamicBroadcast, EapolDynamicBoth}

type Network struct {
	ssid     string          `json:"ssid"`
	scanSsid ScanSSID        `json:"scanSsid,omitempty"`
	bssid    string          `json:"bssid,omitempty"`
	priority uint            `json:"priority,omitempty"`
	mode     Mode            `json:"mode,omitempty"`
	proto    []Proto         `json:"proto,omitempty"`
	keyMngnt []KeyManagement `json:"key_mgmt"`
	authAlg  []AuthAlg       `json:"auth_Alg,omitempty"`
	pairWise []PairWise      `json:"pairwise,omitempty"`
	group    []Group         `json:"group,omitempty"`
	psk      string          `json:"psk,omitempty"`
	eaPol    EapolFlag       `json:"eapol_flags,omitempty"`
	eap      []eapMethod
}

type networkBuilder interface {
	WithSSID(ssid string) networkBuilder
	WithScanSSID(sssid ScanSSID) networkBuilder
	WithBSSID(bssid string) networkBuilder
	WithPriority(prio uint) networkBuilder
	WithMode(mode Mode) networkBuilder
	WithProto(proto ...Proto) networkBuilder
	WithKeyManagement(keyMng ...KeyManagement) networkBuilder
	WithAuthAlg(alg ...AuthAlg) networkBuilder
	WithPairWise(wise ...PairWise) networkBuilder
	WithGroup(group ...Group) networkBuilder
	WithPSK(psk string) networkBuilder
	WithEapolFlag(flag EapolFlag) networkBuilder
	WithEAPMethods(eapMethod ...eapMethod) networkBuilder
	Build() (*Network, error)
}

type NetworkBuilder struct {
	// Ssid Service Set IDentifier
	ssid       string          `json:"ssid"`
	scanSsid   ScanSSID        `json:"scanSsid,omitempty"`
	// Bssid Basic Service Set IDentifier
	bssid      string          `json:"bssid,omitempty"`
	priority   uint            `json:"priority,omitempty"`
	mode       Mode            `json:"mode,omitempty"`
	proto      []Proto         `json:"proto,omitempty"`
	keyMngnt   []KeyManagement `json:"key_mgmt"`
	authAlg    []AuthAlg       `json:"auth_Alg,omitempty"`
	pairWise   []PairWise      `json:"pairwise,omitempty"`
	group      []Group         `json:"group,omitempty"`
	psk        string          `json:"psk,omitempty"`
	eaPol      EapolFlag       `json:"eapol_flags,omitempty"`
	eapMethods []eapMethod
}

func NewNetworkBuilder() networkBuilder {
	netBuilder := NetworkBuilder{
		scanSsid: -1,
		priority: 0,
		mode:     -1,
		eaPol:    -1,
	}
	return &netBuilder
}

// WithSSID sets network name (as announced by the access point).
// An ASCII or hex string enclosed in quotation marks
func (b *NetworkBuilder) WithSSID(ssid string) networkBuilder {
	b.ssid = ssid
	return b
}

// WithScanSSID scan technique; 0 (default) or 1.
// Technique 0 scans for the SSID using a broadcast Probe Request frame while 1 uses a directed Probe Request frame.
// Access points that cloak themselves by not broadcasting their SSID require technique 1, but beware
// that this scheme can cause scanning to take longer to complete.
func (b *NetworkBuilder) WithScanSSID(sssid ScanSSID) networkBuilder {
	b.scanSsid = sssid
	return b
}

// WithBSSID sets network BSSID (typically the MAC address of the access point)
func (b *NetworkBuilder) WithBSSID(bssid string) networkBuilder {
	b.bssid = bssid
	return b
}

// WithPriority of a network when selecting among multiple networks;
// a higher value means a network is more desirable.
// By default networks have priority 0.
// When multiple networks with the same priority are considered for selection,  other information such as security
// policy and signal strength are used to select one
func (b *NetworkBuilder) WithPriority(prio uint) networkBuilder {
	b.priority = prio
	return b
}

// WithMode IEEE 802.11 operation mode; either 0 (infrastructure, default) or 1 (IBSS).
// Note that IBSS (adhoc) mode can only be used with key_mgmt set to NONE (plaintext and static WEP).
func (b *NetworkBuilder) WithMode(mode Mode) networkBuilder {
	b.mode = mode
	return b
}

// WithProto List of acceptable protocols; one or more of: WPA (IEEE 802.11i/D3.0) and RSN (IEEE 802.11i).
// WPA2 is another name for RSN . If not set this defaults to “WPA RSN”.
func (b *NetworkBuilder) WithProto(proto ...Proto) networkBuilder {
	b.proto = make([]Proto, 0)
	b.proto = append(b.proto, proto...)
	return b
}

// WithKeyManagement List of acceptable key management protocols; one or more of: WPA-PSK (WPA pre-shared key),
// WPA-EAP (WPA using EAP authentication),
// IEEE8021X (IEEE 802.1x using EAP authentication and, optionally, dynamically generated WEP keys),
// NONE (plaintext or static WEP keys). If not set this defaults to “WPA-PSK WPA-EAP”.
func (b *NetworkBuilder) WithKeyManagement(keyMng ...KeyManagement) networkBuilder {
	b.keyMngnt = make([]KeyManagement, 0)
	b.keyMngnt = append(b.keyMngnt, keyMng...)
	return b
}

// WithAuthAlg List of allowed IEEE 802.11 authentication algorithms; one or more of:
// OPEN (Open System authentication, required for WPA/WPA2);
// SHARED (Shared Key authentication);
// LEAP (LEAP/Network EAP);
// If not set automatic selection is used (Open System with LEAP enabled if LEAP is allowed as one of the EAP methods).
func (b *NetworkBuilder) WithAuthAlg(alg ...AuthAlg) networkBuilder {
	b.authAlg = make([]AuthAlg, 0)
	b.authAlg = append(b.authAlg, alg...)
	return b
}

// WithPairWise List of acceptable pairwise (unicast) ciphers for WPA; one or more of:
// CCMP (AES in Counter mode with CBC-MAC, RFC 3610, IEEE 802.11i/D7.0);
// TKIP (Temporal Key Integrity Protocol, IEEE 802.11i/D7.0);
// NONE (deprecated);
// If not set this defaults to “CCMP TKIP”.
func (b *NetworkBuilder) WithPairWise(wise ...PairWise) networkBuilder {
	b.pairWise = make([]PairWise, 0)
	b.pairWise = append(b.pairWise, wise...)
	return b
}

// WithGroup List of acceptable group (multicast) ciphers for WPA; one or more of:
// CCMP (AES in Counter mode with CBC-MAC,0 RFC 3610, IEEE 802.11i/D7.0);
// TKIP (Temporal Key Integrity Protocol, IEEE 802.11i/D7.0);
// WEP104 (WEP with 104-bit key), WEP40 (WEP with 40-bit key);
// If not set this defaults to “CCMP TKIP WEP104 WEP40”.
func (b *NetworkBuilder) WithGroup(group ...Group) networkBuilder {
	b.group = make([]Group, 0)
	b.group = append(b.group, group...)
	return b
}

// WithPSK WPA preshared key used in WPA-PSK mode.
// The key is specified as 64 hex digits or as an 8-63 character ASCII passphrase.
// ASCII passphrases are converted to a 256-bit key using the network SSID by the wpa_passphrase(8) utility.
func (b *NetworkBuilder) WithPSK(psk string) networkBuilder {
	b.psk = psk
	return b
}

// WithEapolFlag Dynamic WEP key usage for non-WPA mode, specified as a bit field. Bit 0 (1) forces dynamically
// generated unicast WEP keys to be used. Bit 1 (2) forces dynamically generated broadcast WEP keys to be used.
// By default this is set to 3 (use both)
func (b *NetworkBuilder) WithEapolFlag(flag EapolFlag) networkBuilder {
	b.eaPol = flag
	return b
}

// WithEAPMethods List of acceptable EAP methods; one or more of:
// MD5 (EAP-MD5, cannot be used with WPA, used only as a Phase 2 method with EAP-PEAP or EAP-TTLS);
// MSCHAPV2 (EAP-MSCHAPV2, cannot be used with WPA; used only as a Phase 2 method with EAP-PEAP or EAP-TTLS);
// OTP (EAP-OTP, cannot be used with WPA; used only as a Phase 2 method with EAP-PEAP or EAP-TTLS);
// GTC (EAP-GTC, cannot be used with WPA; used only as a Phase 2 method with EAP-PEAP or EAP-TTLS);
// TLS (EAP-TLS, client and server certificate), PEAP (EAP-PEAP, with tunneled EAP authentication);
// TTLS (EAP-TTLS, with tunneled EAP or PAP/CHAP/MSCHAP/MSCHAPV2 authentication);
// If not set this defaults to all available methods compiled in to wpa_supplicant(8);
// Note that by default wpa_supplicant(8) is compiled with EAP support.
func (b *NetworkBuilder) WithEAPMethods(eapMethods ...eapMethod) networkBuilder {
	b.eapMethods = make([]eapMethod, 0)
	b.eapMethods = append(b.eapMethods, eapMethods...)
	return b
}

func (b *NetworkBuilder) Build() (*Network, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}
	netConfig := Network{
		ssid:     b.ssid,
		scanSsid: b.scanSsid,
		bssid:    b.bssid,
		priority: b.priority,
		mode:     b.mode,
		proto:    b.proto,
		keyMngnt: b.keyMngnt,
		authAlg:  b.authAlg,
		pairWise: b.pairWise,
		group:    b.group,
		psk:      b.psk,
		eaPol:    b.eaPol,
		eap:      b.eapMethods,
	}
	return &netConfig, nil
}

func (b *NetworkBuilder) validate() error {
	if b.ssid == "" {
		if !contains(b.keyMngnt, IEEE8021X) {
			return errors.New("no ssid specified and no IEEE8021X key mngt. specify at least one")
		}
	}
	if b.scanSsid != -1 && !contains(scanSlice, b.scanSsid) {
		return errors.New("invalid value for scanSSID")
	}
	if b.mode != -1 && !contains(modeSlice, b.mode) {
		return errors.New("invalid value for mode")
	}
	if !contains(protoSlice, b.proto...) {
		return errors.New("invalid value for proto")
	}
	if !contains(keyMngtSlice, b.keyMngnt...) {
		return errors.New("invalid value for key management")
	}
	if !contains(authAlgSlice, b.authAlg...) {
		return errors.New("invalid value for auth alg")
	}
	if !contains(pairWiseSlice, b.pairWise...) {
		return errors.New("invalid value for pair wise")
	}
	if !contains(groupSlice, b.group...) {
		return errors.New("invalid value for group")
	}
	if b.eaPol != -1 && !contains(eapFlagSlice, b.eaPol) {
		return errors.New("invalid value for eapol flag")
	}
	if len(b.eapMethods) == 0 {
		return errors.New("at least one eap method must be specifed")
	}
	return nil
}

func contains[T comparable](slice []T, values ...T) bool {
	for i := 0; i < len(values); i++ {
		if !slices.Contains(slice, values[i]) {
			return false
		}
	}
	return true
}

func (net *Network) ToConfigString() string {
	builder := strings.Builder{}
	builder.WriteString("network={\n")
	if net.ssid != "" {
		builder.WriteString(fmt.Sprintf("ssid=\"%s\"\n", net.ssid))
	}
	if net.scanSsid != -1 {
		builder.WriteString(fmt.Sprintf("scan_ssid=%d\n", net.scanSsid))
	}
	if net.bssid != "" {
		builder.WriteString(fmt.Sprintf("bssid=\"%s\"\n", net.bssid))
	}
	if net.priority != 0 {
		builder.WriteString(fmt.Sprintf("priority=%d\n", net.priority))
	}
	if net.mode != -1 {
		builder.WriteString(fmt.Sprintf("mode=%d\n", net.mode))
	}
	if net.proto != nil && len(net.proto) > 0 {
		builder.WriteString(fmt.Sprintf("proto=%s", net.proto[0]))
		for i := 1; i < len(net.proto); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.proto[i]))
		}
		builder.WriteString("\n")
	}
	if net.keyMngnt != nil && len(net.keyMngnt) > 0 {
		builder.WriteString(fmt.Sprintf("key_mgmt=%s", net.keyMngnt[0]))
		for i := 1; i < len(net.keyMngnt); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.keyMngnt[i]))
		}
		builder.WriteString("\n")
	}
	if net.authAlg != nil && len(net.authAlg) > 0 {
		builder.WriteString(fmt.Sprintf("auth_alg=%s", net.authAlg[0]))
		for i := 1; i < len(net.authAlg); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.authAlg[i]))
		}
		builder.WriteString("\n")
	}
	if net.pairWise != nil && len(net.pairWise) > 0 {
		builder.WriteString(fmt.Sprintf("pairwise=%s", net.pairWise[0]))
		for i := 1; i < len(net.pairWise); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.pairWise[i]))
		}
		builder.WriteString("\n")
	}
	if net.group != nil && len(net.group) > 0 {
		builder.WriteString(fmt.Sprintf("group=%s", net.group[0]))
		for i := 1; i < len(net.group); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.group[i]))
		}
		builder.WriteString("\n")
	}
	if net.psk != "" {
		builder.WriteString(fmt.Sprintf("ssid=%s\n", net.psk))
	}
	if net.eaPol != -1 {
		builder.WriteString(fmt.Sprintf("eapol_flags=%d\n", net.eaPol))
	}
	if net.eap != nil && len(net.eap) > 0 {
		builder.WriteString(fmt.Sprintf("eap=%s", net.eap[0].GetEAPName()))
		for i := 1; i < len(net.eap); i++ {
			builder.WriteString(fmt.Sprintf(" %s", net.eap[i].GetEAPName()))
		}
		builder.WriteString("\n")
		for i := 0; i < len(net.eap); i++ {
			builder.WriteString(net.eap[i].ToConfigString())
		}
	}
	builder.WriteString("}\n")
	return builder.String()
}

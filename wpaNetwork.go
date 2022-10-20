package wpaSuppDBusLib

import (
	"errors"
	"golang.org/x/exp/slices"
)

type ScanSSID byte
type Mode byte
type Proto string
type AuthAlg string
type KeyManagement string
type PairWise string
type Group string
type EapolFlag uint8

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
	eap      []eapProvider
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
	WithEAPMethods(eapMethod ...eapProvider) networkBuilder
	Build() (*Network, error)
}

type NetworkBuilder struct {
	// Ssid Service Set IDentifier
	Ssid     string   `json:"ssid"`
	ScanSsid ScanSSID `json:"scanSsid,omitempty"`
	// Bssid Basic Service Set IDentifier
	Bssid      string          `json:"bssid,omitempty"`
	Priority   uint            `json:"priority,omitempty"`
	Mode       Mode            `json:"mode,omitempty"`
	Proto      []Proto         `json:"proto,omitempty"`
	KeyMngnt   []KeyManagement `json:"key_mgmt"`
	AuthAlg    []AuthAlg       `json:"auth_Alg,omitempty"`
	PairWise   []PairWise      `json:"pairwise,omitempty"`
	Group      []Group         `json:"group,omitempty"`
	Psk        string          `json:"psk,omitempty"`
	EaPol      EapolFlag       `json:"eapol_flags,omitempty"`
	EapMethods []eapProvider
}

func NewNetworkBuilder() networkBuilder {
	netBuilder := NetworkBuilder{}
	return &netBuilder
}

// WithSSID sets network name (as announced by the access point).
// An ASCII or hex string enclosed in quotation marks
func (b *NetworkBuilder) WithSSID(ssid string) networkBuilder {
	b.Ssid = ssid
	return b
}

// WithScanSSID scan technique; 0 (default) or 1.
// Technique 0 scans for the SSID using a broadcast Probe Request frame while 1 uses a directed Probe Request frame.
// Access points that cloak themselves by not broadcasting their SSID require technique 1, but beware
// that this scheme can cause scanning to take longer to complete.
func (b *NetworkBuilder) WithScanSSID(sssid ScanSSID) networkBuilder {
	b.ScanSsid = sssid
	return b
}

// WithBSSID sets network BSSID (typically the MAC address of the access point)
func (b *NetworkBuilder) WithBSSID(bssid string) networkBuilder {
	b.Bssid = bssid
	return b
}

// WithPriority of a network when selecting among multiple networks;
// a higher value means a network is more desirable.
// By default networks have priority 0.
// When multiple networks with the same priority are considered for selection,  other information such as security
// policy and signal strength are used to select one
func (b *NetworkBuilder) WithPriority(prio uint) networkBuilder {
	b.Priority = prio
	return b
}

// WithMode IEEE 802.11 operation mode; either 0 (infrastructure, default) or 1 (IBSS).
// Note that IBSS (adhoc) mode can only be used with key_mgmt set to NONE (plaintext and static WEP).
func (b *NetworkBuilder) WithMode(mode Mode) networkBuilder {
	b.Mode = mode
	return b
}

// WithProto List of acceptable protocols; one or more of: WPA (IEEE 802.11i/D3.0) and RSN (IEEE 802.11i).
// WPA2 is another name for RSN . If not set this defaults to “WPA RSN”.
func (b *NetworkBuilder) WithProto(proto ...Proto) networkBuilder {
	b.Proto = make([]Proto, 0)
	b.Proto = append(b.Proto, proto...)
	return b
}

// WithKeyManagement List of acceptable key management protocols; one or more of: WPA-PSK (WPA pre-shared key),
// WPA-EAP (WPA using EAP authentication),
// IEEE8021X (IEEE 802.1x using EAP authentication and, optionally, dynamically generated WEP keys),
// NONE (plaintext or static WEP keys). If not set this defaults to “WPA-PSK WPA-EAP”.
func (b *NetworkBuilder) WithKeyManagement(keyMng ...KeyManagement) networkBuilder {
	b.KeyMngnt = make([]KeyManagement, 0)
	b.KeyMngnt = append(b.KeyMngnt, keyMng...)
	return b
}

// WithAuthAlg List of allowed IEEE 802.11 authentication algorithms; one or more of:
// OPEN (Open System authentication, required for WPA/WPA2);
// SHARED (Shared Key authentication);
// LEAP (LEAP/Network EAP);
// If not set automatic selection is used (Open System with LEAP enabled if LEAP is allowed as one of the EAP methods).
func (b *NetworkBuilder) WithAuthAlg(alg ...AuthAlg) networkBuilder {
	b.AuthAlg = make([]AuthAlg, 0)
	b.AuthAlg = append(b.AuthAlg, alg...)
	return b
}

// WithPairWise List of acceptable pairwise (unicast) ciphers for WPA; one or more of:
// CCMP (AES in Counter mode with CBC-MAC, RFC 3610, IEEE 802.11i/D7.0);
// TKIP (Temporal Key Integrity Protocol, IEEE 802.11i/D7.0);
// NONE (deprecated);
// If not set this defaults to “CCMP TKIP”.
func (b *NetworkBuilder) WithPairWise(wise ...PairWise) networkBuilder {
	b.PairWise = make([]PairWise, 0)
	b.PairWise = append(b.PairWise, wise...)
	return b
}

// WithGroup List of acceptable group (multicast) ciphers for WPA; one or more of:
// CCMP (AES in Counter mode with CBC-MAC,0 RFC 3610, IEEE 802.11i/D7.0);
// TKIP (Temporal Key Integrity Protocol, IEEE 802.11i/D7.0);
// WEP104 (WEP with 104-bit key), WEP40 (WEP with 40-bit key);
// If not set this defaults to “CCMP TKIP WEP104 WEP40”.
func (b *NetworkBuilder) WithGroup(group ...Group) networkBuilder {
	b.Group = make([]Group, 0)
	b.Group = append(b.Group, group...)
	return b
}

// WithPSK WPA preshared key used in WPA-PSK mode.
// The key is specified as 64 hex digits or as an 8-63 character ASCII passphrase.
// ASCII passphrases are converted to a 256-bit key using the network SSID by the wpa_passphrase(8) utility.
func (b *NetworkBuilder) WithPSK(psk string) networkBuilder {
	b.Psk = psk
	return b
}

// WithEapolFlag Dynamic WEP key usage for non-WPA mode, specified as a bit field. Bit 0 (1) forces dynamically
// generated unicast WEP keys to be used. Bit 1 (2) forces dynamically generated broadcast WEP keys to be used.
// By default this is set to 3 (use both)
func (b *NetworkBuilder) WithEapolFlag(flag EapolFlag) networkBuilder {
	b.EaPol = flag
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
func (b *NetworkBuilder) WithEAPMethods(eapMethod ...eapProvider) networkBuilder {
	b.EapMethods = make([]eapProvider, 0)
	b.EapMethods = append(b.EapMethods, eapMethod...)
	return b
}

func (b *NetworkBuilder) Build() (*Network, error) {
	err := b.validate()
	if err != nil {
		return nil, err
	}
	netConfig := Network{
		ssid:     b.Ssid,
		scanSsid: b.ScanSsid,
		bssid:    b.Bssid,
		priority: b.Priority,
		mode:     b.Mode,
		proto:    b.Proto,
		keyMngnt: b.KeyMngnt,
		authAlg:  b.AuthAlg,
		pairWise: b.PairWise,
		group:    b.Group,
		psk:      b.Psk,
		eaPol:    b.EaPol,
		eap:      b.EapMethods,
	}
	return &netConfig, nil
}

func (b *NetworkBuilder) validate() error {
	if b.Ssid == "" {
		if !contains(b.KeyMngnt, IEEE8021X) {
			return errors.New("no ssid specified and no IEEE8021X key mngt. specify at least one")
		}
	}
	if !contains(scanSlice, b.ScanSsid) {
		return errors.New("invalid value for scanSSID")
	}
	if !contains(modeSlice, b.Mode) {
		return errors.New("invalid value for mode")
	}
	if !contains(protoSlice, b.Proto...) {
		return errors.New("invalid value for proto")
	}
	if !contains(keyMngtSlice, b.KeyMngnt...) {
		return errors.New("invalid value for key management")
	}
	if !contains(authAlgSlice, b.AuthAlg...) {
		return errors.New("invalid value for auth alg")
	}
	if !contains(pairWiseSlice, b.PairWise...) {
		return errors.New("invalid value for pair wise")
	}
	if !contains(groupSlice, b.Group...) {
		return errors.New("invalid value for group")
	}
	if !contains(eapFlagSlice, b.EaPol) {
		return errors.New("invalid value for eapol flag")
	}
	if len(b.EapMethods) == 0 {
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

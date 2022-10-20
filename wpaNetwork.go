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
	Ssid       string          `json:"ssid"`
	ScanSsid   ScanSSID        `json:"scanSsid,omitempty"`
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

func (b *NetworkBuilder) WithSSID(ssid string) networkBuilder {
	b.Ssid = ssid
	return b
}

func (b *NetworkBuilder) WithScanSSID(sssid ScanSSID) networkBuilder {
	b.ScanSsid = sssid
	return b
}

func (b *NetworkBuilder) WithBSSID(bssid string) networkBuilder {
	b.Bssid = bssid
	return b
}

func (b *NetworkBuilder) WithPriority(prio uint) networkBuilder {
	b.Priority = prio
	return b
}

func (b *NetworkBuilder) WithMode(mode Mode) networkBuilder {
	b.Mode = mode
	return b
}

func (b *NetworkBuilder) WithProto(proto ...Proto) networkBuilder {
	b.Proto = make([]Proto, 0)
	b.Proto = append(b.Proto, proto...)
	return b
}

func (b *NetworkBuilder) WithKeyManagement(keyMng ...KeyManagement) networkBuilder {
	b.KeyMngnt = make([]KeyManagement, 0)
	b.KeyMngnt = append(b.KeyMngnt, keyMng...)
	return b
}

func (b *NetworkBuilder) WithAuthAlg(alg ...AuthAlg) networkBuilder {
	b.AuthAlg = make([]AuthAlg, 0)
	b.AuthAlg = append(b.AuthAlg, alg...)
	return b
}

func (b *NetworkBuilder) WithPairWise(wise ...PairWise) networkBuilder {
	b.PairWise = make([]PairWise, 0)
	b.PairWise = append(b.PairWise, wise...)
	return b
}

func (b *NetworkBuilder) WithGroup(group ...Group) networkBuilder {
	b.Group = make([]Group, 0)
	b.Group = append(b.Group, group...)
	return b
}

func (b *NetworkBuilder) WithPSK(psk string) networkBuilder {
	b.Psk = psk
	return b
}

func (b *NetworkBuilder) WithEapolFlag(flag EapolFlag) networkBuilder {
	b.EaPol = flag
	return b
}

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

func contains[T any](slice []T, values ...T) bool {
	for i := 0; i < len(values); i++ {
		if !slices.Contains(slice, values[i]) {
			return false
		}
	}
	return true
}

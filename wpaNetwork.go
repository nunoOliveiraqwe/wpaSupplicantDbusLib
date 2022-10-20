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
	ssid       string          `json:"ssid"`
	scanSsid   ScanSSID        `json:"scanSsid,omitempty"`
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

func (b *NetworkBuilder) WithSSID(ssid string) networkBuilder {
	b.ssid = ssid
	return b
}

func (b *NetworkBuilder) WithScanSSID(sssid ScanSSID) networkBuilder {
	b.scanSsid = sssid
	return b
}

func (b *NetworkBuilder) WithBSSID(bssid string) networkBuilder {
	b.bssid = bssid
	return b
}

func (b *NetworkBuilder) WithPriority(prio uint) networkBuilder {
	b.priority = prio
	return b
}

func (b *NetworkBuilder) WithMode(mode Mode) networkBuilder {
	b.mode = mode
	return b
}

func (b *NetworkBuilder) WithProto(proto ...Proto) networkBuilder {
	b.proto = make([]Proto, 0)
	b.proto = append(b.proto, proto...)
	return b
}

func (b *NetworkBuilder) WithKeyManagement(keyMng ...KeyManagement) networkBuilder {
	b.keyMngnt = make([]KeyManagement, 0)
	b.keyMngnt = append(b.keyMngnt, keyMng...)
	return b
}

func (b *NetworkBuilder) WithAuthAlg(alg ...AuthAlg) networkBuilder {
	b.authAlg = make([]AuthAlg, 0)
	b.authAlg = append(b.authAlg, alg...)
	return b
}

func (b *NetworkBuilder) WithPairWise(wise ...PairWise) networkBuilder {
	b.pairWise = make([]PairWise, 0)
	b.pairWise = append(b.pairWise, wise...)
	return b
}

func (b *NetworkBuilder) WithGroup(group ...Group) networkBuilder {
	b.group = make([]Group, 0)
	b.group = append(b.group, group...)
	return b
}

func (b *NetworkBuilder) WithPSK(psk string) networkBuilder {
	b.psk = psk
	return b
}

func (b *NetworkBuilder) WithEapolFlag(flag EapolFlag) networkBuilder {
	b.eaPol = flag
	return b
}

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

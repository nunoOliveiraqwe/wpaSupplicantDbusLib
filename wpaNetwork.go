package wpaSuppDBusLib

type KeyManagement string
type ScanSSID byte
type Mode byte
type Proto string
type AuthAlg string
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
	WithEAPMethods(eapMethod ...eapProvider) networkBuilder
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
	b.KeyMngnt = append(b.KeyMngnt, keyMng...)
	return b
}

func (b *NetworkBuilder) WithPairWise(wise ...PairWise) networkBuilder {
	//TODO implement me
	panic("implement me")
}

func (b *NetworkBuilder) WithGroup(group ...Group) networkBuilder {
	//TODO implement me
	panic("implement me")
}

func (b *NetworkBuilder) WithPSK(psk string) networkBuilder {
	//TODO implement me
	panic("implement me")
}

func (b *NetworkBuilder) WithEAPMethods(eapMethod ...eapProvider) networkBuilder {
	//TODO implement me
	panic("implement me")
}

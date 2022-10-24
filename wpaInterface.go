package wpaSuppDBusLib

import (
	"errors"
	"fmt"
	"strings"
)

type EapolVersion uint8
type ApScan uint8
type FastReauth byte

const (
	EapolV1       EapolVersion = 1
	EapolV2       EapolVersion = 2
	ApScanOff     ApScan       = 0
	ApScanOn      ApScan       = 1
	ApScanOnV2    ApScan       = 2
	FastReauthOn  FastReauth   = 1
	FastReauthOff FastReauth   = 0
)

var eapolVersionSlice = []EapolVersion{EapolV1, EapolV2}
var apScanSlice = []ApScan{ApScanOn, ApScanOnV2, ApScanOff}
var fastReauthSlice = []FastReauth{FastReauthOn, FastReauthOff}

var defaultCtrInterface = "/var/run/wpa_supplicant"
var defaultCtrInterfaceGroup = ""
var defaultEapolVersion = EapolV1
var defaultApScan = ApScanOn
var defaultFastReauth = FastReauthOn

type WPAInterface struct {
	ctrlInterface      string       `json:"ctrl_interface"`
	ctrlInterfaceGroup string       `json:"ctrl_interface_group,omitempty"`
	eapolVersion       EapolVersion `json:"eapol_version,omitempty"`
	apScan             ApScan       `json:"ap_scan,omitempty"`
	fastReauth         FastReauth   `json:"fast_reauth,omitempty"`
	network            []Network    `json:"network"`
}

func (wpa *WPAInterface) ToConfigString() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("ctrl_interface=%s\n", wpa.ctrlInterface))
	if wpa.ctrlInterfaceGroup != defaultCtrInterfaceGroup {
		builder.WriteString(fmt.Sprintf("ctrl_interface_group=%s\n", wpa.ctrlInterfaceGroup))
	}
	if wpa.eapolVersion != defaultEapolVersion {
		builder.WriteString(fmt.Sprintf("eapol_version=%d\n", wpa.eapolVersion))
	}
	if wpa.apScan != defaultApScan {
		builder.WriteString(fmt.Sprintf("ap_scan=%d\n", wpa.apScan))
	}
	if wpa.fastReauth != defaultFastReauth {
		builder.WriteString(fmt.Sprintf("fast_reauth=%d\n", wpa.fastReauth))
	}
	if wpa.network != nil && len(wpa.network) > 0 {
		for i := 0; i < len(wpa.network); i++ {
			builder.WriteString(wpa.network[i].ToConfigString())
		}
	}
	return builder.String()
}

type wpaInterfaceBuilder interface {
	WithCtrlInterface(iface string) wpaInterfaceBuilder
	WithCtrlInterfaceGroup(ctrlIfaceGroup string) wpaInterfaceBuilder
	WithEapolVersion(version EapolVersion) wpaInterfaceBuilder
	WithApScan(scan ApScan) wpaInterfaceBuilder
	WithFastReauth(reauth FastReauth) wpaInterfaceBuilder
	WithNetwork(net ...Network) wpaInterfaceBuilder
	Build() (*WPAInterface, error)
}

type WpaInterfaceBuilder struct {
	ctrlInterface      string       `json:"ctrl_interface"`
	ctrlInterfaceGroup string       `json:"ctrl_interface_group,omitempty"`
	eapolVersion       EapolVersion `json:"eapol_version,omitempty"`
	apScan             ApScan       `json:"ap_scan,omitempty"`
	fastReauth         FastReauth   `json:"fast_reauth,omitempty"`
	network            []Network    `json:"network"`
}

func NewWpaInterfaceBuilder() wpaInterfaceBuilder {
	builder := WpaInterfaceBuilder{
		ctrlInterface:      defaultCtrInterface,
		ctrlInterfaceGroup: defaultCtrInterfaceGroup,
		eapolVersion:       defaultEapolVersion,
		apScan:             defaultApScan,
		fastReauth:         defaultFastReauth,
	}
	return &builder
}

func (w *WpaInterfaceBuilder) WithCtrlInterface(iface string) wpaInterfaceBuilder {
	w.ctrlInterface = iface
	return w
}

func (w *WpaInterfaceBuilder) WithCtrlInterfaceGroup(ctrlIfaceGroup string) wpaInterfaceBuilder {
	w.ctrlInterfaceGroup = ctrlIfaceGroup
	return w
}

func (w *WpaInterfaceBuilder) WithEapolVersion(version EapolVersion) wpaInterfaceBuilder {
	w.eapolVersion = version
	return w
}

func (w *WpaInterfaceBuilder) WithApScan(scan ApScan) wpaInterfaceBuilder {
	w.apScan = scan
	return w
}

func (w *WpaInterfaceBuilder) WithFastReauth(reauth FastReauth) wpaInterfaceBuilder {
	w.fastReauth = reauth
	return w
}

func (w *WpaInterfaceBuilder) WithNetwork(net ...Network) wpaInterfaceBuilder {
	w.network = make([]Network, 0)
	w.network = append(w.network, net...)
	return w
}

func (w WpaInterfaceBuilder) Build() (*WPAInterface, error) {
	err := w.validate()
	if err != nil {
		return nil, err
	}
	wpaIf := WPAInterface{
		ctrlInterface:      w.ctrlInterface,
		ctrlInterfaceGroup: w.ctrlInterfaceGroup,
		eapolVersion:       w.eapolVersion,
		apScan:             w.apScan,
		fastReauth:         w.fastReauth,
		network:            w.network,
	}
	return &wpaIf, err
}

func (w WpaInterfaceBuilder) validate() error {
	if w.ctrlInterface == "" {
		return errors.New("invalid value for ctrl interface")
	}
	if !contains(eapolVersionSlice, w.eapolVersion) {
		return errors.New("invalid value for eapol version")
	}
	if !contains(apScanSlice, w.apScan) {
		return errors.New("invalid value for apscan")
	}
	if !contains(fastReauthSlice, w.fastReauth) {
		return errors.New("invalid value for fast reauth")
	}
	if w.network == nil || len(w.network) == 0 {
		return errors.New("no networks configured. at least one network must be provided")
	}
	return nil
}

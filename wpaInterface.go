package wpaSuppDBusLib

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

var defaultCtrInterface = "/var/run/wpa_supplicant"
var defaultCtrInterfaceGroup = ""
var defaultEapolVersion = EapolV1
var defaultApScan = ApScanOn
var defaultFastReauth = FastReauthOn

type NetworkInterface struct {
	ctrlInterface      string       `json:"ctrl_interface"`
	ctrlInterfaceGroup string       `json:"ctrl_interface_group,omitempty"`
	eapolVersion       EapolVersion `json:"eapol_version,omitempty"`
	apScan             ApScan       `json:"ap_scan,omitempty"`
	fastReauth         FastReauth   `json:"fast_reauth,omitempty"`
	network            []Network    `json:"network"`
}

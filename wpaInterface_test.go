package wpaSuppDBusLib

import (
	"strings"
	"testing"
)

func TestWpaInterfaceToConfTextSmallExample(t *testing.T) {
	//expected text is modeled after https://wiki.archlinux.org/title/wpa_supplicant#802.1x/radius
	expectedConfText := "ctrl_interface=/run/wpa_supplicant\nap_scan=0\nnetwork={\n  key_mgmt=IEEE8021X\n  eap=PEAP\n  identity=\"user_name\"\n  password=\"user_password\"\n  phase2=\"auth=MSCHAPV2\"\n}\n"

	//build auth
	peapAuthBuilder := NewPEAPBuilder()
	eapPEAP, _ := peapAuthBuilder.WithIdentity("user_name").WithPassword("user_password").WithInnerAuthType(InnerAuthMsChapV2).Build()

	//build net block
	netBUilder := NewNetworkBuilder()
	network, _ := netBUilder.WithKeyManagement(IEEE8021X).WithEAPMethods(eapPEAP).Build()

	//build interface
	ifBuilder := NewWpaInterfaceBuilder()
	x, _ := ifBuilder.WithApScan(ApScanOff).WithNetwork(*network).WithCtrlInterface("/run/wpa_supplicant").Build()
	confStr := x.ToConfigString()

	if !strings.EqualFold(expectedConfText, confStr) {
		t.Errorf("config strings don't match")
	}
}

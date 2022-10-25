package wpaSuppDBusLib

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"golang.org/x/exp/slices"
	"os"
	"path"
)

type Driver string

/*
drivers:
nl80211 = Linux nl80211/cfg80211
wext = Linux wireless extensions (generic)
wired = Wired Ethernet driver
macsec_linux = MACsec Ethernet driver for Linux
none = no driver (RADIUS server/WPS ER)
*/
const (
	DriverNL80211     Driver = "nl80211"
	DriverWext        Driver = "wext"
	DriverWired       Driver = "wired"
	DriverMacSecLinux Driver = "macsec_linux"
	DriverNone        Driver = "None"
)

type WpaSupplicantDbus struct {
	dbusCon              *dbus.Conn
	logger               Logger
	EapMethods           []string
	WFDIEs               []byte
	Capabilities         []string
	DebugShowKeys        bool
	DebugTimeStamp       bool
	DebugLevel           string
	CreatedWPAInterfaces map[string]WPAInterface
}

func NewWpaSupplicantDaemonWithLogger(logger Logger) (*WpaSupplicantDbus, error) {
	con, err := newConn()
	if err != nil {
		return nil, err
	}
	supDaemon := WpaSupplicantDbus{dbusCon: con, logger: logger, CreatedWPAInterfaces: make(map[string]WPAInterface)}
	return &supDaemon, nil
}

func NewWpaSupplicantDaemon() (*WpaSupplicantDbus, error) {
	logger := newDefaultLogger()
	return NewWpaSupplicantDaemonWithLogger(logger)
}

func (wpaDbus *WpaSupplicantDbus) Close() error {
	return wpaDbus.dbusCon.Close()
}

func (wpaDbus *WpaSupplicantDbus) CreateInterface(interfaceName, bridgeName string, driver Driver, wpaInterface WPAInterface, pathToSaveInterfaceConfig string, signalChannel chan<- *dbus.Signal) (dbus.ObjectPath, error) {
	confStr := wpaInterface.ToConfigString()
	fileName := ""
	if driver == DriverWired {
		fileName = fmt.Sprintf("wpa_supplicant-wired-%s", interfaceName)
	} else {
		fileName = fmt.Sprintf("wpa_supplicant-%s", interfaceName)
	}
	fullPath := path.Join(pathToSaveInterfaceConfig, fileName)
	err := os.WriteFile(fullPath, []byte(confStr), 0600)
	if err != nil {
		return "", err
	}
	ifPath, err := createInterface(wpaDbus, interfaceName, bridgeName, driver, fullPath, signalChannel)
	if err != nil {
		return "", err
	}
	wpaDbus.CreatedWPAInterfaces[string(ifPath)] = wpaInterface
	return ifPath, nil
}

func (wpaDbus *WpaSupplicantDbus) ExpectDisconnect(wpaInterfaceName string) error {
	return expectDisconnect(wpaDbus, wpaInterfaceName)
}

func (wpaDbus *WpaSupplicantDbus) GetInterface(systemNetworkInterfaceName string) (dbus.ObjectPath, error) {
	return getInterface(wpaDbus, systemNetworkInterfaceName)
}

func (wpaDbus *WpaSupplicantDbus) RemoveInterface(wpaInterfaceName dbus.ObjectPath) error {
	return removeInterface(wpaDbus, wpaInterfaceName)
}

func (wpaDbus *WpaSupplicantDbus) ReadAllProperties() error {
	readWFDIEs(wpaDbus)
	readCapabilities(wpaDbus)
	readDebugShowKeys(wpaDbus)
	readDebugTimeStamp(wpaDbus)
	readDebugLevel(wpaDbus)
	readEapMethods(wpaDbus)
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

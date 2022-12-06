package wpaSuppDBusLib

import (
	"fmt"
	"github.com/godbus/dbus/v5"
	"os"
	"path"
	"reflect"
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

func (wpaDbus *WpaSupplicantDbus) CreateInterface(interfaceName, bridgeName string, driver Driver, wpaInterface WPAInterface, pathToSaveInterfaceConfig string, stateChangeChan chan string) (dbus.ObjectPath, error) {
	confStr := wpaInterface.ToConfigString()
	fileName := ""
	if driver == DriverWired {
		fileName = fmt.Sprintf("wpa_supplicant-wired-%s.conf", interfaceName)
	} else {
		fileName = fmt.Sprintf("wpa_supplicant-%s.conf", interfaceName)
	}
	fullPath := path.Join(pathToSaveInterfaceConfig, fileName)
	err := os.WriteFile(fullPath, []byte(confStr), 0600)
	if err != nil {
		return "", err
	}
	ifPath, err := createInterface(wpaDbus, interfaceName, bridgeName, driver, fullPath, stateChangeChan)
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

func contains(slice interface{}, values interface{}) bool {
	if reflect.TypeOf(slice).Kind() == reflect.Slice || reflect.TypeOf(slice).Kind() == reflect.Array {
		list := reflect.ValueOf(slice)
		if reflect.TypeOf(values).Kind() == reflect.Slice || reflect.TypeOf(values).Kind() == reflect.Array {
			valueList := reflect.ValueOf(values)
			for j := 0; j < valueList.Len(); j++ {
				if valueList.Index(j).Interface() == nil {
					continue
				}
				found := false
				for i := 0; i < list.Len(); i++ {
					if valueList.Index(j).Interface() == list.Index(i).Interface() {
						found = true
					}
				}
				if !found {
					return false
				}
			}
		} else {
			found := false
			for i := 0; i < list.Len(); i++ {
				if values == list.Index(i).Interface() {
					found = true
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

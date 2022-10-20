package wpaSuppDBusLib

import (
	"github.com/godbus/dbus/v5"
)

var dbusWPAname = "fi.w1.wpa_supplicant1"

var dbusWPAObjectPath = dbus.ObjectPath("/fi/w1/wpa_supplicant1")

func newConn() (*dbus.Conn, error) {
	con, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	return con, nil
}

func readWFDIEs(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "WFDIEs").Store(&wpaDbus.WFDIEs)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	return nil
}

func readCapabilities(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "Capabilities").Store(&wpaDbus.Capabilities)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	return nil
}

func readDebugShowKeys(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "DebugShowKeys").Store(&wpaDbus.DebugShowKeys)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	return nil
}

func readDebugTimeStamp(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "DebugTimeStamp").Store(&wpaDbus.DebugTimeStamp)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	return nil
}

func readDebugLevel(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "DebugLevel").Store(&wpaDbus.DebugLevel)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	return nil
}

//
//func readEapMethods(wpaDbus *WpaSupplicantDbus) error {
//	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
//	var availableEAPMethods []string
//	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "EapMethods").Store(&availableEAPMethods)
//	if err != nil {
//		wpaDbus.logger.Error(err)
//		return err
//	}
//	wpaDbus.RawEapMethods = availableEAPMethods
//	availableEapProviderKeys := eapRegistryGetMatchingKeys(availableEAPMethods)
//
//}

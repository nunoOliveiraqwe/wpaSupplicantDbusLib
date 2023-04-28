package wpaSuppDBusLib

import (
	"errors"
	"github.com/godbus/dbus/v5"
)

var dbusWPAname = "fi.w1.wpa_supplicant1"

var dbusWPAInterfacename = "fi.w1.wpa_supplicant1.Interface"

var dbusWPAObjectPath = dbus.ObjectPath("/fi/w1/wpa_supplicant1")

func newConn() (*dbus.Conn, error) {
	con, err := dbus.ConnectSystemBus()
	if err != nil {
		return nil, err
	}
	return con, nil
}

func createInterface(wpaDbus *WpaSupplicantDbus, interfaceName, bridgeName string, driver Driver, pathToSaveInterfaceConfig string, stateChangeChan chan string) (dbus.ObjectPath, error) {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	var result interface{}
	argMap := make(map[string]interface{})
	argMap["Ifname"] = interfaceName
	argMap["BridgeIfname"] = bridgeName
	argMap["Driver"] = string(driver)
	argMap["ConfigFile"] = pathToSaveInterfaceConfig
	err := obj.Call(dbusWPAname+".CreateInterface", 0, argMap).Store(&result)
	if err != nil {
		wpaDbus.logger.Error(err)
		return "", err
	}
	if _, ok := result.(dbus.ObjectPath); !ok {
		return "", errors.New("unknown return type from dbus. expected string")
	}
	interfaceNameRet := result.(dbus.ObjectPath)
	objInterface := wpaDbus.dbusCon.Object(dbusWPAInterfacename, interfaceNameRet)

	if err = wpaDbus.dbusCon.AddMatchSignal(
		dbus.WithMatchObjectPath(objInterface.Path()),
		dbus.WithMatchInterface(objInterface.Destination()),
		dbus.WithMatchMember("PropertiesChanged"),
	); err != nil {
		panic(err)
	}
	signalChan := make(chan *dbus.Signal)
	wpaDbus.dbusCon.Signal(signalChan)

	go stateChangeListenFunc(stateChangeChan, signalChan)

	return interfaceNameRet, nil
}

func stateChangeListenFunc(stateChangeChan chan string, signalChan chan *dbus.Signal) {
	for changedProp := range signalChan {
		for i := 0; i < len(changedProp.Body); i++ {
			noTypeMap := changedProp.Body[i].(map[string]dbus.Variant)
			if value, contains := noTypeMap["State"]; contains {
				stateChangeChan <- value.String()
			}
		}
	}
}

func removeInterface(wpaDbus *WpaSupplicantDbus, wpaInterfaceName dbus.ObjectPath) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call(dbusWPAname+".RemoveInterface", 0, wpaInterfaceName).Err
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	delete(wpaDbus.CreatedWPAInterfaces, string(wpaInterfaceName))
	return nil
}

func expectDisconnect(wpaDbus *WpaSupplicantDbus, wpaInterfaceName string) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	err := obj.Call(dbusWPAname+".RemoveInterface", 0, dbus.ObjectPath(wpaInterfaceName)).Err
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	delete(wpaDbus.CreatedWPAInterfaces, wpaInterfaceName)
	return nil
}

func getInterface(wpaDbus *WpaSupplicantDbus, networkInterfaceName string) (dbus.ObjectPath, error) {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	var result interface{}
	err := obj.Call(dbusWPAname+".GetInterface", 0, networkInterfaceName).Store(&result)
	if err != nil {
		wpaDbus.logger.Error(err)
		return "", err
	}
	if _, ok := result.(dbus.ObjectPath); !ok {
		return "", errors.New("unknown return type from dbus. expected string")
	}
	return result.(dbus.ObjectPath), nil
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

func readEapMethods(wpaDbus *WpaSupplicantDbus) error {
	obj := wpaDbus.dbusCon.Object(dbusWPAname, dbusWPAObjectPath)
	var availableEAPMethods []string
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, dbusWPAname, "EapMethods").Store(&availableEAPMethods)
	if err != nil {
		wpaDbus.logger.Error(err)
		return err
	}
	wpaDbus.EapMethods = availableEAPMethods
	return nil
}

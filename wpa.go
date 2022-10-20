package wpaSuppDBusLib

import (
	"github.com/godbus/dbus/v5"
)

type WpaSupplicantDbus struct {
	dbusCon                *dbus.Conn
	logger                 Logger
	EapMethods             []string
	WFDIEs                 []byte
	Capabilities           []string
	DebugShowKeys          bool
	DebugTimeStamp         bool
	DebugLevel             string
	availableProvidersKeys []string
}

func NewWpaSupplicantDaemonWithLogger(logger Logger) (*WpaSupplicantDbus, error) {
	con, err := newConn()
	if err != nil {
		return nil, err
	}
	supDaemon := WpaSupplicantDbus{dbusCon: con, logger: logger}
	return &supDaemon, nil
}

func NewWpaSupplicantDaemon() (*WpaSupplicantDbus, error) {
	logger := newDefaultLogger()
	return NewWpaSupplicantDaemonWithLogger(logger)
}

func (wpaDbus *WpaSupplicantDbus) Close() error {
	return wpaDbus.dbusCon.Close()
}

func (wpaDbus *WpaSupplicantDbus) CreateInterface() {

}

func (wpaDbus *WpaSupplicantDbus) ExpectDisconnect() {

}

func (wpaDbus *WpaSupplicantDbus) GetInterface() {

}

func (wpaDbus *WpaSupplicantDbus) RemoveInterface() {

}

func (wpaDbus *WpaSupplicantDbus) ReadAllProperties() error {
	readWFDIEs(wpaDbus)
	readCapabilities(wpaDbus)
	readDebugShowKeys(wpaDbus)
	readDebugTimeStamp(wpaDbus)
	readDebugLevel(wpaDbus)
	return nil
}

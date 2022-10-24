package main

import (
	"flag"
	"fmt"
	wpaSuppDBusLib "git.dev.zgrp.net/litecom/libs/wpaSupplicantDbusLib"
	"github.com/godbus/dbus/v5"
	"log"
)

func main() {
	caCertPath := flag.String("caCertPath", "", "path to ca cert")
	identity := flag.String("identity", "", "tls identity")
	clientCertPath := flag.String("clientCertPath", "", "path to client cert")
	privateKeyPath := flag.String("privateKeyPath", "", "path to private key")
	privateKeyPassword := flag.String("privateKeyPassword", "", "private key password")
	interfaceName := flag.String("interfaceName", "eth0", "network interface name")
	wpaCtrlInterface := flag.String("wpaCtrlInterface", "/run/wpa_supplicant", "path to control interface of wpa supplicant")
	storagePathToWpaConfFiles := flag.String("storagePathToWpaConfFile", "/etc/wpa_supplicant/", "path where to store wpa interface config file")
	flag.Parse()

	tlsAuthBuilder := wpaSuppDBusLib.NewTLSBuilder()
	tlsEap, _ := tlsAuthBuilder.
		WithIdentity(*identity).
		WithCaCertPath(*caCertPath).
		WithClientCertPath(*clientCertPath).
		WithPrivateKeyPath(*privateKeyPath).
		WithPrivateKeyPassword(*privateKeyPassword).Build()

	//build net block
	netBUilder := wpaSuppDBusLib.NewNetworkBuilder()
	network, _ := netBUilder.WithKeyManagement(wpaSuppDBusLib.IEEE8021X).WithEapolFlag(wpaSuppDBusLib.EapolOff).WithEAPMethods(tlsEap).Build()

	//build interface
	ifBuilder := wpaSuppDBusLib.NewWpaInterfaceBuilder()
	x, _ := ifBuilder.WithApScan(wpaSuppDBusLib.ApScanOff).WithNetwork(*network).WithCtrlInterface(*wpaCtrlInterface).Build()
	confStr := x.ToConfigString()
	log.Println(confStr)

	supDaemon, _ := wpaSuppDBusLib.NewWpaSupplicantDaemon()
	signalChannel := make(chan *dbus.Signal, 10)
	supDaemon.CreateInterface(*interfaceName, "", wpaSuppDBusLib.DriverWired, *x, *storagePathToWpaConfFiles, signalChannel)
	for msg := range signalChannel {
		fmt.Println(msg)
	}
}

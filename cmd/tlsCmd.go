package main

import (
	"flag"
	wpaSuppDBusLib "git.dev.zgrp.net/litecom/libs/wpaSupplicantDbusLib"
	"log"
)

func main() {
	caCertPath := flag.String("caCertPath", "", "path to ca cert")
	identity := flag.String("identity", "", "tls identity")
	clientCertPath := flag.String("clientCertPath", "", "path to client cert")
	privateKeyPath := flag.String("privateKeyPath", "", "path to private key")
	privateKeyPassword := flag.String("privateKeyPassword", "", "private key password")
	interfaceName := flag.String("interfaceName", "enp0s3", "network interface name")
	wpaCtrlInterface := flag.String("wpaCtrlInterface", "/run/wpa_supplicant", "path to control interface of wpa supplicant")
	storagePathToWpaConfFiles := flag.String("storagePathToWpaConfFile", "/etc/wpa_supplicant/", "path where to store wpa interface config file")
	flag.Parse()

	//build auth type TLS
	tlsAuthBuilder := wpaSuppDBusLib.NewTLSBuilder()
	tlsEap, err := tlsAuthBuilder.
		WithIdentity(*identity).
		WithCaCertPath(*caCertPath).
		WithClientCertPath(*clientCertPath).
		WithPrivateKeyPath(*privateKeyPath).
		WithPrivateKeyPassword(*privateKeyPassword).Build()
	if err != nil {
		log.Fatalln(err)
	}

	//build a wired network with prev auth type
	netBuilder := wpaSuppDBusLib.NewNetworkBuilder()
	network, _ := netBuilder.WithKeyManagement(wpaSuppDBusLib.IEEE8021X).WithEapolFlag(wpaSuppDBusLib.EapolOff).WithEAPMethods(tlsEap).Build()

	//build the WPA interface definition
	ifBuilder := wpaSuppDBusLib.NewWpaInterfaceBuilder()
	wpaInterface, err := ifBuilder.WithApScan(wpaSuppDBusLib.ApScanOff).WithNetwork(*network).WithCtrlInterface(*wpaCtrlInterface).Build()
	if err != nil {
		log.Fatalln(err)
	}
	//print config
	confStr := wpaInterface.ToConfigString()
	log.Println(confStr)

	//use supplicant API to set it all up
	supplicantAPI, _ := wpaSuppDBusLib.NewWpaSupplicantAPI()
	stateChannel := make(chan string, 10)
	dbusPath, err := supplicantAPI.CreateInterface(*interfaceName, "", wpaSuppDBusLib.DriverWired, *wpaInterface, *storagePathToWpaConfFiles, stateChannel)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("new interface is " + dbusPath)
	for message := range stateChannel {
		log.Println("Con state -> " + message)
	}
}

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
	stateChannel := make(chan string, 10)
	_, err := supDaemon.CreateInterface(*interfaceName, "", wpaSuppDBusLib.DriverWired, *x, *storagePathToWpaConfFiles, stateChannel)
	if err != nil {
		log.Fatalln(err)
	}
	for message := range stateChannel {
		log.Println("Con state -> " + message)
	}
}

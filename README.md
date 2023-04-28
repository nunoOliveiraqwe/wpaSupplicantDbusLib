# Go WPA Supplicant DBus Library

This is a Go library that provides an abstraction over the details of interacting with the WPA Supplicant via the DBus interface. 
It enables developers to easily build  WPA Supplicant configurations, including support for various authentication methods such as TLS, TTLS, EAP, and MD5.

## Installation

You can install the library with the following command:

```
go get github.com/nunooliveiraqwe/wpa_supplicant_dbus_lib
```

## Usage

The example code below demonstrates how to use the library to configure the WPA Supplicant for a wired configuration with TLS auth:

```go
package main

import (
	"flag"
	"log"

	"github.com/your_username/wpa_supplicant_dbus_lib"
)

func main() {
	// Parse command line arguments
	caCertPath := flag.String("caCertPath", "", "path to ca cert")
	identity := flag.String("identity", "", "tls identity")
	clientCertPath := flag.String("clientCertPath", "", "path to client cert")
	privateKeyPath := flag.String("privateKeyPath", "", "path to private key")
	privateKeyPassword := flag.String("privateKeyPassword", "", "private key password")
	interfaceName := flag.String("interfaceName", "enp0s3", "network interface name")
	wpaCtrlInterface := flag.String("wpaCtrlInterface", "/run/wpa_supplicant", "path to control interface of wpa supplicant")
	storagePathToWpaConfFiles := flag.String("storagePathToWpaConfFile", "/etc/wpa_supplicant/", "path where to store wpa interface config file")
	flag.Parse()

	// Build auth type TLS
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

	// Build a wired network with previous auth type
	netBuilder := wpaSuppDBusLib.NewNetworkBuilder()
	network, _ := netBuilder.WithKeyManagement(wpaSuppDBusLib.IEEE8021X).WithEapolFlag(wpaSuppDBusLib.EapolOff).WithEAPMethods(tlsEap).Build()

	// Build the WPA interface definition
	ifBuilder := wpaSuppDBusLib.NewWpaInterfaceBuilder()
	wpaInterface, err := ifBuilder.WithApScan(wpaSuppDBusLib.ApScanOff).WithNetwork(*network).WithCtrlInterface(*wpaCtrlInterface).Build()
	if err != nil {
		log.Fatalln(err)
	}

	// Print config
	confStr := wpaInterface.ToConfigString()
	log.Println(confStr)

	// Use supplicant API to set it all up
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
``
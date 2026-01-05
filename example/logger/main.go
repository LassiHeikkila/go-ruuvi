package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"

	"github.com/LassiHeikkila/go-ruuvi/ruuvi"
)

var jsonOnlyMode = flag.Bool("jsonOnly", false, "Only print Ruuvi data as JSON, nothing else. Useful for piping output to jq or storing in a file, etc.")

func println(a ...interface{}) {
	if !*jsonOnlyMode {
		fmt.Println(a...)
	}
}

func printf(format string, a ...interface{}) {
	if !*jsonOnlyMode {
		fmt.Printf(format, a...)
	}
}

func onStateChanged(d gatt.Device, s gatt.State) {
	println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		println("scanning...")
		d.Scan([]gatt.UUID{}, true) // report duplicates since we want to keep receiving adverts from sensors
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if !ruuvi.IsAdvertisementFromRuuviTag(a.ManufacturerData) {
		return
	}
	printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	println("  Local Name        =", a.LocalName)
	println("  TX Power Level    =", a.TxPowerLevel)
	println("  Manufacturer Data =", hex.EncodeToString(a.ManufacturerData))
	println("  Service Data      =", a.ServiceData)

	outputRuuviDinternal/pkg/rawv1/dataformat.goata(a.ManufacturerData)
}

func main() {
	// disable default log output, some noise is coming from gatt library
	log.SetOutput(ioutil.Discard)
	flag.Parse()
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		printf("Failed to open device, err: %s\n", err)
		os.Exit(1)
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)
	select {}
}

func outputRuuviData(b []byte) {
	if len(b) == 0 {
		return
	}
	println("Device has advertisement payload:", hex.EncodeToString(b))

	advert, err := ruuvi.ProcessAdvertisement(b)
	if err != nil {
		println("Error processing bytes:", err)
		return
	}

	advertJson, err := advert.MarshalJSON()
	if err != nil {
		println("Error marshalling ruuvi data")
		return
	}
	fmt.Println(string(advertJson))
}

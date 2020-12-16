package main

import (
	"encoding/hex"
	"flag"
	"fmt"

	"tinygo.org/x/bluetooth"

	"github.com/LassiHeikkila/go-ruuvi/ruuvi"
)

var adapter = bluetooth.DefaultAdapter

var jsonOnlyMode = flag.Bool("jsonOnly", false, "Only print Ruuvi data as JSON, nothing else. Useful for piping output to jq or storing in a file, etc.")

func main() {
	flag.Parse()
	// Enable BLE interface.
	err := adapter.Enable()
	if err != nil {
		panic("Failed to enable BLE stack: " + err.Error())
	}

	// Start scanning.
	Print("Scanning...")
	err = adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
		Print("Found device:", device.Address.String(), device.RSSI, device.LocalName())
		b := device.Bytes()
		if len(b) == 0 {
			return
		}
		Print("Device has advertisement payload:", hex.EncodeToString(b))

		advert, err := ruuvi.ProcessAdvertisement(b)
		if err != nil {
			Print("Error processing bytes:", err)
		}

		advertJson, err := advert.MarshalJSON()
		if err != nil {
			Print("Error marshalling ruuvi data")
		}
		fmt.Println(string(advertJson))
	})
	if err != nil {
		panic("Failed to start scan: " + err.Error())
	}
}

func Print(a ...interface{}) {
	if !*jsonOnlyMode {
		fmt.Println(a...)
	}
}

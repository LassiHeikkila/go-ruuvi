package ruuvi

import (
	"fmt"

	"github.com/LassiHeikkila/go-ruuvi/internal/pkg/rawv1"
	"github.com/LassiHeikkila/go-ruuvi/internal/pkg/rawv2"
)

// AdvertisementData is an interface abstracting away raw data from Ruuvitag BLE advertisements
//
// It provides methods to easily get the interesting values without having to manually parse byte arrays
type AdvertisementData interface {
	// DataFormat returns format of underlying data
	DataFormat() int8

	// Temperature returns measured temperature in degrees Celsius
	Temperature() (float64, error)

	// Humidity returns measured humidity as percentage
	Humidity() (float64, error)

	// Pressure returns measured atmospheric pressure with unit Pa (pascal)
	Pressure() (int, error)

	// AccelerationX returns the acceleration in X axis with unit G, if supported by data format
	AccelerationX() (float64, error)

	// AccelerationY returns the acceleration in Y axis with unit G, if supported by data format
	AccelerationY() (float64, error)

	// AccelerationZ returns the acceleration in Z axis with unit G, if supported by data format
	AccelerationZ() (float64, error)

	// BatteryVoltage returns battery voltage with unit V (volt), if supported by data format
	BatteryVoltage() (float64, error)

	// TransmissionPower returns transmission power with unit dBm, if supported by data format
	TransmissionPower() (float64, error)

	// MovementCounter returns number of movements detected by accelerometer, if supported by data format
	MovementCounter() (int, error)

	// MeasurementSequenceNumber returns measurement sequence number, if supported by data format
	MeasurementSequenceNumber() (int, error)

	// MACAddress returns MAC address of broadcasting ruuvitag, if supported by data format
	MACAddress() ([]byte, error)

	// RawData returns the raw bytes. Make sure to copy the data, or it may be overwritten by the next broadcast.
	RawData() []byte

	// Copy copies the raw bytes internally so the AdvertisementData object is safe to use for a longer time.
	// Without Copy(), incoming BLE packets can overwrite the bytes
	Copy()
}

// ProcessAdvertisement processes the given bytes and returns AdvertisementData or error
// error will be nil if AdvertisementData is valid (given data was valid and of a supported format)
// error will be non-nil if given data was invalid or of an unsupported format.
func ProcessAdvertisement(data []byte) (AdvertisementData, error) {
	switch data[0] {
	case 0x3:
		return rawv1.NewDataRAWv1(data)
	case 0x5:
		return rawv2.NewDataRAWv2(data)
	}
	return nil, newUnsupportedData("Package does not support this data format (yet)")
}

// UnsupportedData is an error returned when package does not know how to handle given data
type UnsupportedData struct {
	description string
}

func newUnsupportedData(desc string) *UnsupportedData {
	return &UnsupportedData{description: desc}
}

func (ud *UnsupportedData) Error() string {
	return fmt.Sprint("Unsupported data:", ud.description)
}

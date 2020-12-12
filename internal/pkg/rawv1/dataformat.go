package rawv1

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// DataRAWv1 is a concrete implementation of AdvertisementData interface
// Data format is described here: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-3-rawv1
type DataRAWv1 struct {
	rawBytes []byte
}

// NewDataRAWv1 returns pointer to DataRAWv1 wrapping
func NewDataRAWv1(d []byte) (*DataRAWv1, error) {
	if determineDataVersion(d) != 3 {
		return nil, errors.New("Data is not RAWv1 (3)")
	}
	if len(d) < 14 {
		return nil, errors.New("Data is too short to be valid, expected 14 bytes")
	}

	return &DataRAWv1{rawBytes: d}, nil
}

func determineDataVersion(d []byte) int8 {
	return int8(d[0])
}

func dataNotAvailable(whatData string) error {
	return fmt.Errorf("%s is not available with data format RAWv1 (3)", whatData)
}

// Copy copies the raw bytes internally so the AdvertisementData object is safe to use for a longer time.
// Without Copy(), incoming BLE packets can overwrite the bytes
func (d *DataRAWv1) Copy() {
	c := make([]byte, len(d.rawBytes))
	copy(c[:], d.rawBytes[:])

	d.rawBytes = c
}

// DataFormat returns format of underlying data
func (d *DataRAWv1) DataFormat() int8 { return 3 }

// Temperature returns measured temperature in degrees Celsius
func (d *DataRAWv1) Temperature() (float64, error) {
	t1 := d.rawBytes[2] & 0b01111111
	negative := d.rawBytes[2]&0b10000000 > 0
	t2 := d.rawBytes[3]

	if t2 > 99 {
		return 0, errors.New("Temperature fractional part exceeds maximum value")
	}
	var mult float64
	if negative {
		mult = -1.0
	} else {
		mult = 1.0
	}

	temp := (float64(t1) + (float64(t2) / 100.0)) * mult

	return temp, nil
}

// Humidity returns measured humidity as percentage
func (d *DataRAWv1) Humidity() (float64, error) {
	b := d.rawBytes[1]

	humidity := float64(b) * 0.5
	return humidity, nil
}

// Pressure returns measured atmospheric pressure with unit Pa (pascal)
func (d *DataRAWv1) Pressure() (int, error) {
	pb := d.rawBytes[4:6]

	pres := binary.BigEndian.Uint16(pb)
	return int(pres) + 50000, nil
}

// AccelerationX returns the acceleration in X axis with unit G, if supported by data format
func (d *DataRAWv1) AccelerationX() (float64, error) {
	b := d.rawBytes[6:8]
	acc := int16(binary.BigEndian.Uint16(b))
	gs := float64(acc) / 1000.0
	return gs, nil
}

// AccelerationY returns the acceleration in Y axis with unit G, if supported by data format
func (d *DataRAWv1) AccelerationY() (float64, error) {
	b := d.rawBytes[8:10]
	acc := int16(binary.BigEndian.Uint16(b))
	gs := float64(acc) / 1000.0
	return gs, nil
}

// AccelerationZ returns the acceleration in Z axis with unit G, if supported by data format
func (d *DataRAWv1) AccelerationZ() (float64, error) {
	b := d.rawBytes[10:12]
	acc := int16(binary.BigEndian.Uint16(b))
	gs := float64(acc) / 1000.0
	return gs, nil
}

// BatteryVoltage returns battery voltage with unit V (volt), if supported by data format
func (d *DataRAWv1) BatteryVoltage() (float64, error) {
	b := d.rawBytes[12:14]
	millivolts := binary.BigEndian.Uint16(b)
	return float64(millivolts) / 1000, nil
}

// TransmissionPower returns transmission power with unit dBm, if supported by data format
func (d *DataRAWv1) TransmissionPower() (float64, error) {
	return 0, dataNotAvailable("TX power")
}

// MovementCounter returns number of movements detected by accelerometer, if supported by data format
func (d *DataRAWv1) MovementCounter() (int, error) {
	return 0, dataNotAvailable("Movement counter")
}

// MeasurementSequenceNumber returns measurement sequence number, if supported by data format
func (d *DataRAWv1) MeasurementSequenceNumber() (int, error) {
	return 0, dataNotAvailable("Measurement sequence number")
}

// MACAddress returns MAC address (48 bits / 6 bytes) of broadcasting ruuvitag, if supported by data format
func (d *DataRAWv1) MACAddress() ([]byte, error) {
	return nil, dataNotAvailable("MAC address")
}

// RawData returns the raw bytes. Make sure to copy the data, or it may be overwritten by the next broadcast.
func (d *DataRAWv1) RawData() []byte {
	return d.rawBytes
}

package rawv2

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// DataRAWv2 is a concrete implementation of AdvertisementData interface
// Data format is described here: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-5-rawv2
type DataRAWv2 struct {
	rawBytes []byte
}

// InvalidValue is error returned when raw data contains data specified as invalid,
// i.e. 0xFFFF for unsigned values or 0x8000 for signed values
type InvalidValue struct {
	whatMeasurement string
}

func newInvalidValue(whatMeasurement string) *InvalidValue {
	return &InvalidValue{whatMeasurement: whatMeasurement}
}

func (iv *InvalidValue) Error() string {
	return fmt.Sprintf("Data for %s is invalid", iv.whatMeasurement)
}

// Is makes it possible to use errors.Is() on this error type
func (iv *InvalidValue) Is(target error) bool {
	switch target.(type) {
	case *InvalidValue:
		return true
	default:
		return false
	}
}

// NewDataRAWv2 returns pointer to DataRAWv2 wrapping
func NewDataRAWv2(d []byte) (*DataRAWv2, error) {
	if determineDataVersion(d) != 5 {
		return nil, errors.New("Data is not RAWv2 (5)")
	}
	if len(d) < 24 {
		return nil, errors.New("Data is too short to be valid, expected 14 bytes")
	}

	return &DataRAWv2{rawBytes: d}, nil
}

func determineDataVersion(d []byte) int8 {
	return int8(d[0])
}

func dataNotAvailable(whatData string) error {
	return fmt.Errorf("%s is not available with data format RAWv2 (5)", whatData)
}

// Copy copies the raw bytes internally so the AdvertisementData object is safe to use for a longer time.
// Without Copy(), incoming BLE packets can overwrite the bytes
func (d *DataRAWv2) Copy() {
	c := make([]byte, len(d.rawBytes))
	copy(c[:], d.rawBytes[:])

	d.rawBytes = c
}

// DataFormat returns format of underlying data
func (d *DataRAWv2) DataFormat() int8 { return 5 }

// Temperature returns measured temperature in degrees Celsius
func (d *DataRAWv2) Temperature() (float64, error) {
	b := d.rawBytes[1:3]

	u := binary.BigEndian.Uint16(b)

	if u == 0x8000 {
		return 0.0, newInvalidValue("temperature")
	}

	temp := float64(int16(u)) * 0.005

	return temp, nil
}

// Humidity returns measured humidity as percentage
func (d *DataRAWv2) Humidity() (float64, error) {
	b := d.rawBytes[3:5]

	v := binary.BigEndian.Uint16(b)
	if v == 0xFFFF {
		return 0.0, newInvalidValue("humidity")
	}
	humidity := float64(v) * 0.0025
	return humidity, nil
}

// Pressure returns measured atmospheric pressure with unit Pa (pascal)
func (d *DataRAWv2) Pressure() (int, error) {
	pb := d.rawBytes[5:7]

	pres := binary.BigEndian.Uint16(pb)
	if pres == 0xFFFF {
		return 0, newInvalidValue("pressure")
	}
	return int(pres) + 50000, nil
}

// AccelerationX returns the acceleration in X axis with unit G, if supported by data format
func (d *DataRAWv2) AccelerationX() (float64, error) {
	b := d.rawBytes[7:9]
	u := binary.BigEndian.Uint16(b)
	if u == 0x8000 {
		return 0.0, newInvalidValue("acceleration-x")
	}
	acc := int16(u)
	gs := float64(acc) / 1000.0
	return gs, nil
}

// AccelerationY returns the acceleration in Y axis with unit G, if supported by data format
func (d *DataRAWv2) AccelerationY() (float64, error) {
	b := d.rawBytes[9:11]
	u := binary.BigEndian.Uint16(b)
	if u == 0x8000 {
		return 0.0, newInvalidValue("acceleration-y")
	}
	acc := int16(u)
	gs := float64(acc) / 1000.0
	return gs, nil
}

// AccelerationZ returns the acceleration in Z axis with unit G, if supported by data format
func (d *DataRAWv2) AccelerationZ() (float64, error) {
	b := d.rawBytes[11:13]
	u := binary.BigEndian.Uint16(b)
	if u == 0x8000 {
		return 0.0, newInvalidValue("acceleration-z")
	}
	acc := int16(u)
	gs := float64(acc) / 1000.0
	return gs, nil
}

// BatteryVoltage returns battery voltage with unit V (volt), if supported by data format
func (d *DataRAWv2) BatteryVoltage() (float64, error) {
	b := d.rawBytes[13:15]
	v := binary.BigEndian.Uint16(b)

	if v == 0xFFFF {
		return 0.0, newInvalidValue("battery voltage")
	}

	v = (v & 0b1111111111100000) >> 5

	return (float64(v) / 1000) + 1.6, nil
}

// TransmissionPower returns transmission power with unit dBm, if supported by data format
func (d *DataRAWv2) TransmissionPower() (float64, error) {
	b := d.rawBytes[13:15]
	v := binary.BigEndian.Uint16(b)
	if v == 0xFFFF {
		return 0.0, newInvalidValue("tx power")
	}

	v = v & 0b0000000000011111

	return (float64(v) * 2) - 40.0, nil
}

// MovementCounter returns number of movements detected by accelerometer, if supported by data format
func (d *DataRAWv2) MovementCounter() (int, error) {
	b := d.rawBytes[15]

	if b == 0xFF {
		return 0, newInvalidValue("movement counter")
	}

	return int(b), nil
}

// MeasurementSequenceNumber returns measurement sequence number, if supported by data format
func (d *DataRAWv2) MeasurementSequenceNumber() (int, error) {
	b := d.rawBytes[16:18]

	v := binary.BigEndian.Uint16(b)

	if v == 0xFFFF {
		return 0, newInvalidValue("measurement sequence number")
	}

	return int(v), nil
}

// MACAddress returns MAC address (48 bits / 6 bytes) of broadcasting ruuvitag, if supported by data format
func (d *DataRAWv2) MACAddress() ([]byte, error) {
	b := d.rawBytes[18:24]

	if bytes.Equal(b, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}) {
		return nil, newInvalidValue("MAC address")
	}

	return b, nil
}

// RawData returns the raw bytes. Make sure to copy the data, or it may be overwritten by the next broadcast.
func (d *DataRAWv2) RawData() []byte {
	return d.rawBytes
}

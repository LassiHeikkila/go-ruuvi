package df6

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
)

type calibration struct {
	calibrationInProgress bool
	vocB9 bool
	noxB9 bool
}

type DataFormat6 struct {
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


func dataNotAvailable(whatData string) error {
	return fmt.Errorf("%s is not available with data format 6", whatData)
}

// NewDataFormat6 returns pointer to DataFormat6
func NewDataFormat6(d []byte) (*DataFormat6, error) {
	if determineDataVersion(d) != 6 {
		return nil, errors.New("Data is not format 6")
	}
	if len(d) < 20 {
		return nil, errors.New("Data is too short to be valid, expected 20 bytes")
	}

	return &DataFormat6{rawBytes: d}, nil
}

func determineDataVersion(d []byte) int8 {
	return int8(d[0])
}

// Implement the interface

// DataFormat returns format of underlying data
func (d *DataFormat6) DataFormat() int8 {
	return 6
}

// Temperature returns measured temperature in degrees Celsius
func (d *DataFormat6) Temperature() (float64, error) {
	b := d.rawBytes[1:3]

	u := binary.BigEndian.Uint16(b)

	if u == 0x8000 {
		return 0.0, newInvalidValue("temperature")
	}

	temp := float64(int16(u)) * 0.005

	return temp, nil
}

// Humidity returns measured humidity as percentage
func (d *DataFormat6) Humidity() (float64, error) {
	b := d.rawBytes[3:5]

	u := binary.BigEndian.Uint16(b)

	if u == 0xFFFF {
		return 0.0, newInvalidValue("humidity")
	}

	humidity := float64(u) * 0.0025

	return humidity, nil
}

// Pressure returns measured atmospheric pressure with unit Pa (pascal)
func (d *DataFormat6) Pressure() (int, error) {
	b := d.rawBytes[5:7]

	u := binary.BigEndian.Uint16(b)

	if u == 0xFFFF {
		return 0, newInvalidValue("pressure")
	}

	pressure := int(u) + 50000

	return pressure, nil
}

// CO2 returns CO2 level in ppm, if supported by data format
func (d *DataFormat6) CO2() (int, error) {
	b := d.rawBytes[9:11]

	u := binary.BigEndian.Uint16(b)

	if u == 0xFFFF {
		return 0, newInvalidValue("CO2")
	}

	return int(u), nil
}

// PM2p5 returns particulate matter 2.5 concentration with unit µg/m³, if supported by data format
func (d *DataFormat6) PM2p5() (float64, error) {
	b := d.rawBytes[7:9]

	u := binary.BigEndian.Uint16(b)

	if u == 0xFFFF {
		return 0.0, newInvalidValue("PM2p5")
	}

	return float64(u) * 0.1, nil
}

// VOC returns volatile organic compounds index, if supported by data format
func (d *DataFormat6) VOC() (int, error) {
	b := d.rawBytes[11]
	f := (d.rawBytes[16] & 0b01000000) >> 6

	u := uint16(b) << 1 | uint16(f)

	if u == 0x1FF {
		return 0, newInvalidValue("VOC")
	}

	return int(u), nil
}

// NOX returns nitrous oxide index, if supported by data format
func (d *DataFormat6) NOX() (int, error) {
	b := d.rawBytes[12]
	f := (d.rawBytes[16] & 0b10000000) >> 7

	u := uint16(b) << 1 | uint16(f)

	if u == 0x1FF {
		return 0, newInvalidValue("NOX")
	}

	return int(u), nil
}

// Luminosity returns luminosity with unit lux, if supported by data format
func (d *DataFormat6) Luminosity() (float64, error) {
	b := d.rawBytes[13]

	const (
		MAX_VALUE = 65535
		MAX_CODE = 254
		DELTA = 0.04366281452346112 // math.Log(MAX_VALUE + 1) / MAX_CODE
	)
	
	if b == 0xFF {
		return 0.0, newInvalidValue("luminosity")
	}

	value := math.Exp(float64(b) * DELTA) - 1

	return value, nil
}

// AccelerationX returns the acceleration in X axis with unit G, if supported by data format
func (d *DataFormat6) AccelerationX() (float64, error) {
	return 0.0, dataNotAvailable("acceleration-x")
}

// AccelerationY returns the acceleration in Y axis with unit G, if supported by data format
func (d *DataFormat6) AccelerationY() (float64, error) {
	return 0.0, dataNotAvailable("acceleration-y")
}

// AccelerationZ returns the acceleration in Z axis with unit G, if supported by data format
func (d *DataFormat6) AccelerationZ() (float64, error) {
	return 0.0, dataNotAvailable("acceleration-z")
}

// BatteryVoltage returns battery voltage with unit V (volt), if supported by data format
func (d *DataFormat6) BatteryVoltage() (float64, error) {
	return 0.0, dataNotAvailable("batteryVoltage")
}

// TransmissionPower returns transmission power with unit dBm, if supported by data format
func (d *DataFormat6) TransmissionPower() (float64, error) {
	return 0.0, dataNotAvailable("transmissionPower")
}

// MovementCounter returns number of movements detected by accelerometer, if supported by data format
func (d *DataFormat6) MovementCounter() (int, error) {
	return 0, dataNotAvailable("movementCounter")
}

// MeasurementSequenceNumber returns measurement sequence number, if supported by data format
func (d *DataFormat6) MeasurementSequenceNumber() (int, error) {
	b := d.rawBytes[15]

	return int(b), nil
}

// MACAddress returns last 3 bytes of MAC address of broadcasting ruuvitag, if supported by data format
func (d *DataFormat6) MACAddress() ([]byte, error) {
	b := d.rawBytes[17:20]

	return b, nil
}

// RawData returns the raw bytes. Make sure to copy the data, or it may be overwritten by the next broadcast.
func (d *DataFormat6) RawData() []byte {
	return d.rawBytes
}

// Copy copies the raw bytes internally so the AdvertisementData object is safe to use for a longer time.
// Without Copy(), incoming BLE packets can overwrite the bytes
func (d *DataFormat6) Copy() {
	c := make([]byte, len(d.rawBytes))
	copy(c[:], d.rawBytes[:])

	d.rawBytes = c
}

// MarshalJSON outputs available data as JSON
func (d *DataFormat6) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, 13)

	m["raw"] = hex.EncodeToString(d.rawBytes)
	m["format"] = d.DataFormat()
	if t, err := d.Temperature(); err == nil {
		m["temperature"] = t
	}
	if h, err := d.Humidity(); err == nil {
		m["humidity"] = h
	}
	if p, err := d.Pressure(); err == nil {
		m["pressure"] = p
	}
	if co2, err := d.CO2(); err == nil {
		m["CO2"] = co2
	}
	if pm2p5, err := d.PM2p5(); err == nil {
		m["PM2p5"] = pm2p5
	}
	if voc, err := d.VOC(); err == nil {
		m["VOC"] = voc
	}
	if nox, err := d.NOX(); err == nil {
		m["NOX"] = nox
	}
	if lum, err := d.Luminosity(); err == nil {
		m["luminosity"] = lum
	}
	if s, err := d.MeasurementSequenceNumber(); err == nil {
		m["meas-seq"] = s
	}
	if mac, err := d.MACAddress(); err == nil && len(mac) == 6 {
		m["mac"] = fmt.Sprintf("%x:%x:%x", mac[0], mac[1], mac[2])
	}
	m["calibrationInProgress"], _ = d.CalibrationInProgress()

	return json.Marshal(&m)
}

func (d *DataFormat6) CalibrationInProgress() (bool, error) {
	b := d.rawBytes[16]
	if (b & 1) == 1 {
		// calibration in progress
		return true, nil
	}
	return false, nil
}
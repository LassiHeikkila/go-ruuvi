package rawv2

import (
	"errors"
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type result struct {
	dataFormat      int
	temperature     float64
	pressure        int
	humidity        float64
	accelerationX   float64
	accelerationY   float64
	accelerationZ   float64
	txPower         float64
	voltage         float64
	movementCounter int
	measSequence    int
	MAC             []byte
}

var float64FuzzyCompOpt = cmp.Comparer(func(x, y float64) bool {
	delta := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if mean == 0 {
		return true
	}
	return (delta / mean) < 0.00001
})

func TestValidData(t *testing.T) {
	// 0x0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F
	validExampleData := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}

	expectedResult := result{
		dataFormat:      5,
		temperature:     24.3,
		pressure:        100044,
		humidity:        53.49,
		accelerationX:   0.004,
		accelerationY:   -0.004,
		accelerationZ:   1.036,
		txPower:         4.0,
		voltage:         2.977,
		movementCounter: 66,
		measSequence:    205,
		MAC:             []byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
	}

	rawv2, err := NewDataRAWv2(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv2.DataFormat() != 5 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv2.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, expectedResult.temperature, float64FuzzyCompOpt) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv2.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, expectedResult.pressure) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv2.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, expectedResult.humidity, float64FuzzyCompOpt) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv2.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, expectedResult.accelerationX, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv2.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, expectedResult.accelerationY, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv2.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, expectedResult.accelerationZ, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv2.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, expectedResult.voltage, float64FuzzyCompOpt) {
		t.Fatal("Wrong voltage returned:", voltage)
	}

	if txPower, err := rawv2.TransmissionPower(); err != nil {
		t.Fatal("TransmissionPower() returned error")
	} else if !cmp.Equal(txPower, expectedResult.txPower, float64FuzzyCompOpt) {
		t.Fatal("Wrong transmission power returned:", txPower)
	}

	if movCounter, err := rawv2.MovementCounter(); err != nil {
		t.Fatal("MovementCounter() returned error")
	} else if !cmp.Equal(movCounter, expectedResult.movementCounter) {
		t.Fatal("Wrong MovementCounter returned:", movCounter)
	}

	if measSeq, err := rawv2.MeasurementSequenceNumber(); err != nil {
		t.Fatal("MeasurementSequenceNumber() returned error")
	} else if !cmp.Equal(measSeq, expectedResult.measSequence) {
		t.Fatal("Wrong MeasurementSequenceNumber returned:", measSeq)
	}

	if mac, err := rawv2.MACAddress(); err != nil {
		t.Fatal("MACAddress() returned error")
	} else if !cmp.Equal(mac, expectedResult.MAC) {
		t.Fatal("Wrong MAC returned:", mac)
	}
}

func TestMaximumValuesData(t *testing.T) {
	// 0x057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F
	validExampleData := []byte{
		0x05, 0x7F, 0xFF, 0xFF, 0xFE, 0xFF, 0xFE, 0x7F,
		0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0xFF, 0xDE, 0xFE,
		0xFF, 0xFE, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}

	expectedResult := result{
		dataFormat:      5,
		temperature:     163.835,
		pressure:        115534,
		humidity:        163.8350,
		accelerationX:   32.767,
		accelerationY:   32.767,
		accelerationZ:   32.767,
		txPower:         20,
		voltage:         3.646,
		movementCounter: 254,
		measSequence:    65534,
		MAC:             []byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
	}

	rawv2, err := NewDataRAWv2(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv2.DataFormat() != 5 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv2.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, expectedResult.temperature, float64FuzzyCompOpt) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv2.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, expectedResult.pressure) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv2.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, expectedResult.humidity, float64FuzzyCompOpt) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv2.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, expectedResult.accelerationX, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv2.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, expectedResult.accelerationY, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv2.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, expectedResult.accelerationZ, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv2.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, expectedResult.voltage, float64FuzzyCompOpt) {
		t.Fatal("Wrong voltage returned:", voltage)
	}

	if txPower, err := rawv2.TransmissionPower(); err != nil {
		t.Fatal("TransmissionPower() returned error")
	} else if !cmp.Equal(txPower, expectedResult.txPower, float64FuzzyCompOpt) {
		t.Fatal("Wrong transmission power returned:", txPower)
	}

	if movCounter, err := rawv2.MovementCounter(); err != nil {
		t.Fatal("MovementCounter() returned error")
	} else if !cmp.Equal(movCounter, expectedResult.movementCounter) {
		t.Fatal("Wrong MovementCounter returned:", movCounter)
	}

	if measSeq, err := rawv2.MeasurementSequenceNumber(); err != nil {
		t.Fatal("MeasurementSequenceNumber() returned error")
	} else if !cmp.Equal(measSeq, expectedResult.measSequence) {
		t.Fatal("Wrong MeasurementSequenceNumber returned:", measSeq)
	}

	if mac, err := rawv2.MACAddress(); err != nil {
		t.Fatal("MACAddress() returned error")
	} else if !cmp.Equal(mac, expectedResult.MAC) {
		t.Fatal("Wrong MAC returned:", mac)
	}
}

func TestMinimumValuesData(t *testing.T) {
	// 0x058001000000008001800180010000000000CBB8334C884F
	validExampleData := []byte{
		0x05, 0x80, 0x01, 0x00, 0x00, 0x00, 0x00, 0x80,
		0x01, 0x80, 0x01, 0x80, 0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}

	expectedResult := result{
		dataFormat:      5,
		temperature:     -163.835,
		pressure:        50000,
		humidity:        0.000,
		accelerationX:   -32.767,
		accelerationY:   -32.767,
		accelerationZ:   -32.767,
		txPower:         -40,
		voltage:         1.600,
		movementCounter: 0,
		measSequence:    0,
		MAC:             []byte{0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F},
	}

	rawv2, err := NewDataRAWv2(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv2.DataFormat() != 5 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv2.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, expectedResult.temperature, float64FuzzyCompOpt) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv2.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, expectedResult.pressure) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv2.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, expectedResult.humidity, float64FuzzyCompOpt) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv2.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, expectedResult.accelerationX, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv2.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, expectedResult.accelerationY, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv2.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, expectedResult.accelerationZ, float64FuzzyCompOpt) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv2.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, expectedResult.voltage, float64FuzzyCompOpt) {
		t.Fatal("Wrong voltage returned:", voltage)
	}

	if txPower, err := rawv2.TransmissionPower(); err != nil {
		t.Fatal("TransmissionPower() returned error")
	} else if !cmp.Equal(txPower, expectedResult.txPower, float64FuzzyCompOpt) {
		t.Fatal("Wrong transmission power returned:", txPower)
	}

	if movCounter, err := rawv2.MovementCounter(); err != nil {
		t.Fatal("MovementCounter() returned error")
	} else if !cmp.Equal(movCounter, expectedResult.movementCounter) {
		t.Fatal("Wrong MovementCounter returned:", movCounter)
	}

	if measSeq, err := rawv2.MeasurementSequenceNumber(); err != nil {
		t.Fatal("MeasurementSequenceNumber() returned error")
	} else if !cmp.Equal(measSeq, expectedResult.measSequence) {
		t.Fatal("Wrong MeasurementSequenceNumber returned:", measSeq)
	}

	if mac, err := rawv2.MACAddress(); err != nil {
		t.Fatal("MACAddress() returned error")
	} else if !cmp.Equal(mac, expectedResult.MAC) {
		t.Fatal("Wrong MAC returned:", mac)
	}
}

func TestInvalidValues(t *testing.T) {
	//0x058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF
	invalidExampleData := []byte{
		0x05, 0x80, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0x80,
		0x00, 0x80, 0x00, 0x80, 0x00, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	}

	rawv2, err := NewDataRAWv2(invalidExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv2.DataFormat() != 5 {
		t.Fatal("Wrong data format returned")
	}

	if _, err := rawv2.Temperature(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from Temperature()")
	}
	if _, err := rawv2.Humidity(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from Humidity()")
	}
	if _, err := rawv2.Pressure(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from Pressure()")
	}
	if _, err := rawv2.AccelerationX(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from AccelerationX()")
	}
	if _, err := rawv2.AccelerationY(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from AccelerationY()")
	}
	if _, err := rawv2.AccelerationZ(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from AccelerationZ()")
	}
	if _, err := rawv2.TransmissionPower(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from TransmissionPower()")
	}
	if _, err := rawv2.BatteryVoltage(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from BatteryVoltage()")
	}
	if _, err := rawv2.MovementCounter(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from MovementCounter()")
	}
	if _, err := rawv2.MeasurementSequenceNumber(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from MeasurementSequenceNumber()")
	}
	if _, err := rawv2.MACAddress(); !errors.Is(err, &InvalidValue{}) {
		t.Error("No InvalidValue returned from MACAddress()")
	}
}

func TestRawData(t *testing.T) {
	validExampleData := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	rawv2, _ := NewDataRAWv2(validExampleData)

	if !cmp.Equal(validExampleData, rawv2.RawData()) {
		t.Fatal("RawData() returned different bytes than what was put in!")
	}
}

func TestDataModifiedWithoutCopy(t *testing.T) {
	data := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	rawv2, _ := NewDataRAWv2(data)

	data[2] = 0x00
	data[5] = 0xFF

	b := rawv2.RawData()
	if b[2] != 0x00 {
		t.Fatal("underlying data not modified without Copy()")
	}
	if b[5] != 0xFF {
		t.Fatal("underlying data not modified without Copy()")
	}
}

func TestDataNotModifiedWithCopy(t *testing.T) {
	data := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	rawv2, _ := NewDataRAWv2(data)
	rawv2.Copy()

	data[2] = 0x00
	data[5] = 0xFF

	b := rawv2.RawData()
	if b[2] != 0xFC {
		t.Fatal("underlying data modified after calling Copy()")
	}
	if b[5] != 0xC3 {
		t.Fatal("underlying data modified after calling Copy()")
	}
}

func TestErrorReturnedOnBadInput(t *testing.T) {
	wrongDataFormat := []byte{
		0x02, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB, 0xB8, 0x33, 0x4C, 0x88, 0x4F,
	}
	tooShort := []byte{
		0x05, 0x12, 0xFC, 0x53, 0x94, 0xC3, 0x7C, 0x00,
		0x04, 0xFF, 0xFC, 0x04, 0x0C, 0xAC, 0x36, 0x42,
		0x00, 0xCD, 0xCB,
	}

	if _, err := NewDataRAWv2(wrongDataFormat); err == nil {
		t.Fatal("No error from wrong data format")
	}
	if _, err := NewDataRAWv2(tooShort); err == nil {
		t.Fatal("No error from too short data")
	}
}

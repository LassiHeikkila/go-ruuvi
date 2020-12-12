package rawv1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidData(t *testing.T) {
	validExampleData := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}

	rawv1, err := NewDataRAWv1(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv1.DataFormat() != 3 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv1.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, 26.3) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv1.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, 102766) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv1.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, 20.5) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv1.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, -1.000) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv1.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, -1.726) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv1.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, 0.714) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv1.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, 2.899) {
		t.Fatal("Wrong voltage returned:", voltage)
	}
}

func TestMaximumValuesData(t *testing.T) {
	validExampleData := []byte{
		0x03, 0xFF, 0x7F, 0x63, 0xFF, 0xFF, 0x7F,
		0xFF, 0x7F, 0xFF, 0x7F, 0xFF, 0xFF, 0xFF,
	}

	rawv1, err := NewDataRAWv1(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv1.DataFormat() != 3 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv1.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, 127.99) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv1.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, 115535) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv1.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, 127.5) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv1.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, 32.767) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv1.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, 32.767) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv1.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, 32.767) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv1.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, 65.535) {
		t.Fatal("Wrong voltage returned:", voltage)
	}
}

func TestMinimumValuesData(t *testing.T) {
	validExampleData := []byte{
		0x03, 0x00, 0xFF, 0x63, 0x00, 0x00, 0x80,
		0x01, 0x80, 0x01, 0x80, 0x01, 0x00, 0x00,
	}

	rawv1, err := NewDataRAWv1(validExampleData)
	if err != nil {
		t.Fatal("Error: ", err)
	}

	if rawv1.DataFormat() != 3 {
		t.Fatal("Wrong data format returned")
	}

	if temp, err := rawv1.Temperature(); err != nil {
		t.Fatal("Temperature() returned error")
	} else if !cmp.Equal(temp, -127.99) {
		t.Fatal("Wrong temperature returned:", temp)
	}

	if pres, err := rawv1.Pressure(); err != nil {
		t.Fatal("Pressure() returned error")
	} else if !cmp.Equal(pres, 50000) {
		t.Fatal("Wrong pressure returned:", pres)
	}

	if humid, err := rawv1.Humidity(); err != nil {
		t.Fatal("Humidity() returned error")
	} else if !cmp.Equal(humid, 0.0) {
		t.Fatal("Wrong humidity returned:", humid)
	}

	if accelX, err := rawv1.AccelerationX(); err != nil {
		t.Fatal("AccelerationX() returned error")
	} else if !cmp.Equal(accelX, -32.767) {
		t.Fatal("Wrong AccelerationX returned:", accelX)
	}

	if accelY, err := rawv1.AccelerationY(); err != nil {
		t.Fatal("AccelerationY() returned error")
	} else if !cmp.Equal(accelY, -32.767) {
		t.Fatal("Wrong AccelerationY returned:", accelY)
	}

	if accelZ, err := rawv1.AccelerationZ(); err != nil {
		t.Fatal("AccelerationZ() returned error")
	} else if !cmp.Equal(accelZ, -32.767) {
		t.Fatal("Wrong AccelerationZ returned:", accelZ)
	}

	if voltage, err := rawv1.BatteryVoltage(); err != nil {
		t.Fatal("BatteryVoltage() returned error")
	} else if !cmp.Equal(voltage, 0.000) {
		t.Fatal("Wrong voltage returned:", voltage)
	}
}

func TestErrorReturnedOnUnsupportedValues(t *testing.T) {
	validExampleData := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	rawv1, _ := NewDataRAWv1(validExampleData)

	if _, err := rawv1.TransmissionPower(); err == nil {
		t.Fatal("TransmissionPower() did not return error")
	}

	if _, err := rawv1.MovementCounter(); err == nil {
		t.Fatal("MovementCounter() did not return error")
	}

	if _, err := rawv1.MeasurementSequenceNumber(); err == nil {
		t.Fatal("MeasurementSequenceNumber() did not return error")
	}

	if _, err := rawv1.MACAddress(); err == nil {
		t.Fatal("MACAddress() did not return error")
	}

}

func TestRawData(t *testing.T) {
	validExampleData := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	rawv1, _ := NewDataRAWv1(validExampleData)

	if !cmp.Equal(validExampleData, rawv1.RawData()) {
		t.Fatal("RawData() returned different bytes than what was put in!")
	}
}

func TestDataModifiedWithoutCopy(t *testing.T) {
	data := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	rawv1, _ := NewDataRAWv1(data)

	data[2] = 0x00
	data[5] = 0xFF

	b := rawv1.RawData()
	if b[2] != 0x00 {
		t.Fatal("underlying data not modified without Copy()")
	}
	if b[5] != 0xFF {
		t.Fatal("underlying data not modified without Copy()")
	}
}

func TestDataNotModifiedWithCopy(t *testing.T) {
	data := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	rawv1, _ := NewDataRAWv1(data)
	rawv1.Copy()

	data[2] = 0x00
	data[5] = 0xFF

	b := rawv1.RawData()
	if b[2] != 0x1A {
		t.Fatal("underlying data modified after calling Copy()")
	}
	if b[5] != 0x1E {
		t.Fatal("underlying data modified after calling Copy()")
	}
}

func TestErrorReturnedOnBadInput(t *testing.T) {
	wrongDataFormat := []byte{
		0x02, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	tooShort := []byte{
		0x03, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02,
	}

	if _, err := NewDataRAWv1(wrongDataFormat); err == nil {
		t.Fatal("No error from wrong data format")
	}
	if _, err := NewDataRAWv1(tooShort); err == nil {
		t.Fatal("No error from too short data")
	}
}

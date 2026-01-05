package df6

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Test cases use data from: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-6#test-vectors

func TestValidData(t *testing.T) {
	// Test vector: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-6#case-valid-data
	// 0x06170C5668C79E007000C90501D9XXCD004C884F
	data := []byte{
		0x06, 0x17, 0x0C, 0x56,
		0x68, 0xC7, 0x9E, 0x00,
		0x70, 0x00, 0xC9, 0x05,
		0x01, 0xD9, 0x7F, 0xCD,
		0x00, 0x4C, 0x88, 0x4F,
	}

	const (
		expectDataFormat = 6
		expectTemperature = 29.5
		expectPressure = 101102
		expectHumidity = 55.3
		expectPM2p5 = 11.2
		expectCO2 = 201
		expectVOC = 10
		expectNOX = 2
		expectLuminosity = 13026.67
		expectSequenceNumber = 205
		expectFlagCalibrationInProgress = false
		expectFlagVOCb9 = false
		expectFlagNOXb9 = false
	)
	expectMAC := []byte{0x4C, 0x88, 0x4F}

	d, err := NewDataFormat6(data)
	if err != nil {
		t.Fatalf("unexpected error constructing DataFormat6: %e", err)
	}

	gotDataFormat := d.DataFormat()
	if expectDataFormat != gotDataFormat {
		t.Fatalf("expected DataFormat %d, got %d", expectDataFormat, gotDataFormat)
	}
	
	abs := func(x float64) float64 { if x < 0 { return -x }; return x }

	if got, _ := d.Temperature(); abs(got-expectTemperature) > 0.01 {
		t.Fatalf("expected Temperature %.2f, got %.2f", expectTemperature, got)
	}

	if got, _ := d.Pressure(); int(got) != expectPressure {
		t.Fatalf("expected Pressure %d, got %d", expectPressure, int(got))
	}

	if got, _ := d.Humidity(); abs(got-expectHumidity) > 0.01 {
		t.Fatalf("expected Humidity %.2f, got %.2f", expectHumidity, got)
	}

	if got, _ := d.PM2p5(); abs(got-expectPM2p5) > 0.01 {
		t.Fatalf("expected PM2.5 %.2f, got %.2f", expectPM2p5, got)
	}

	if got, _ := d.CO2(); int(got) != expectCO2 {
		t.Fatalf("expected CO2 %d, got %d", expectCO2, int(got))
	}

	if got, _ := d.VOC(); int(got) != expectVOC {
		t.Fatalf("expected VOC %d, got %d", expectVOC, int(got))
	}

	if got, _ := d.NOX(); int(got) != expectNOX {
		t.Fatalf("expected NOX %d, got %d", expectNOX, int(got))
	}

	if got, _ := d.Luminosity(); abs(got-expectLuminosity) > 0.1 {
		t.Fatalf("expected Luminosity %.2f, got %.2f", expectLuminosity, got)
	}

	if got, _ := d.MeasurementSequenceNumber(); int(got) != expectSequenceNumber {
		t.Fatalf("expected SequenceNumber %d, got %d", expectSequenceNumber, int(got))
	}

	if got, _ := d.CalibrationInProgress(); got != expectFlagCalibrationInProgress {
	 	t.Fatalf("expected FlagCalibrationInProgress %v, got %v", expectFlagCalibrationInProgress, got)
	}

	if got, _ := d.MACAddress(); !bytes.Equal(got, expectMAC) {
		t.Fatalf("expected MAC 0x%X, got 0x%X", expectMAC, got)
	}
}

func TestMaxData(t *testing.T) {
	// Test vector: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-6#case-maximum-values
	// 0x067FFF9C40FFFE27109C40FAFAFEXXFF074C8F4F
	data := []byte{
		0x06, 0x7F, 0xFF, 0x9C,
		0x40, 0xFF, 0xFE, 0x27,
		0x10, 0x9C, 0x40, 0xFA,
		0xFA, 0xFE, 0xFF, 0xFF,
		0x07, 0x4C, 0x8F, 0x4F,
	}
	
	const (
		expectDataFormat = 6
		expectTemperature = 163.835
		expectPressure = 115534
		expectHumidity = 100.000
		expectPM2p5 = 1000.0
		expectCO2 = 40000
		expectVOC = 500
		expectNOX = 500
		expectLuminosity = 65535.00
		expectSequenceNumber = 255
		expectFlagCalibrationInProgress = true
		expectFlagVOCb9 = false
		expectFlagNOXb9 = false
	)
	expectMAC := []byte{0x4C, 0x8F, 0x4F}

	d, err := NewDataFormat6(data)
	if err != nil {
		t.Fatalf("unexpected error constructing DataFormat6: %e", err)
	}

	gotDataFormat := d.DataFormat()
	if expectDataFormat != gotDataFormat {
		t.Fatalf("expected DataFormat %d, got %d", expectDataFormat, gotDataFormat)
	}
	
	abs := func(x float64) float64 { if x < 0 { return -x }; return x }

	if got, _ := d.Temperature(); abs(got-expectTemperature) > 0.01 {
		t.Fatalf("expected Temperature %.2f, got %.2f", expectTemperature, got)
	}

	if got, _ := d.Pressure(); int(got) != expectPressure {
		t.Fatalf("expected Pressure %d, got %d", expectPressure, int(got))
	}

	if got, _ := d.Humidity(); abs(got-expectHumidity) > 0.01 {
		t.Fatalf("expected Humidity %.2f, got %.2f", expectHumidity, got)
	}

	if got, _ := d.PM2p5(); abs(got-expectPM2p5) > 0.01 {
		t.Fatalf("expected PM2.5 %.2f, got %.2f", expectPM2p5, got)
	}

	if got, _ := d.CO2(); int(got) != expectCO2 {
		t.Fatalf("expected CO2 %d, got %d", expectCO2, int(got))
	}

	if got, _ := d.VOC(); int(got) != expectVOC {
		t.Fatalf("expected VOC %d, got %d", expectVOC, int(got))
	}

	if got, _ := d.NOX(); int(got) != expectNOX {
		t.Fatalf("expected NOX %d, got %d", expectNOX, int(got))
	}

	if got, _ := d.Luminosity(); abs(got-expectLuminosity) > 0.1 {
		t.Fatalf("expected Luminosity %.2f, got %.2f", expectLuminosity, got)
	}

	if got, _ := d.MeasurementSequenceNumber(); int(got) != expectSequenceNumber {
		t.Fatalf("expected SequenceNumber %d, got %d", expectSequenceNumber, int(got))
	}

	if got, _ := d.CalibrationInProgress(); got != expectFlagCalibrationInProgress {
	 	t.Fatalf("expected FlagCalibrationInProgress %v, got %v", expectFlagCalibrationInProgress, got)
	}

	if got, _ := d.MACAddress(); !bytes.Equal(got, expectMAC) {
		t.Fatalf("expected MAC 0x%X, got 0x%X", expectMAC, got)
	}
}

func TestMinData(t *testing.T) {
	// Test vector: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-6#case-minimum-values
	// 0x0680010000000000000000000000XX00004C884F
	data := []byte{
		0x06, 0x80, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x4C, 0x88, 0x4F,
	}

	const (
		expectDataFormat = 6
		expectTemperature = -163.835
		expectPressure = 50000
		expectHumidity = 0.0
		expectPM2p5 = 0.0
		expectCO2 = 0
		expectVOC = 0
		expectNOX = 0
		expectLuminosity = 0.0
		expectSequenceNumber = 0
		expectFlagCalibrationInProgress = false
		expectFlagVOCb9 = false
		expectFlagNOXb9 = false
	)
	expectMAC := []byte{0x4C, 0x88, 0x4F}

	d, err := NewDataFormat6(data)
	if err != nil {
		t.Fatalf("unexpected error constructing DataFormat6: %e", err)
	}

	gotDataFormat := d.DataFormat()
	if expectDataFormat != gotDataFormat {
		t.Fatalf("expected DataFormat %d, got %d", expectDataFormat, gotDataFormat)
	}
	
	abs := func(x float64) float64 { if x < 0 { return -x }; return x }

	if got, _ := d.Temperature(); abs(got-expectTemperature) > 0.01 {
		t.Fatalf("expected Temperature %.2f, got %.2f", expectTemperature, got)
	}

	if got, _ := d.Pressure(); int(got) != expectPressure {
		t.Fatalf("expected Pressure %d, got %d", expectPressure, int(got))
	}

	if got, _ := d.Humidity(); abs(got-expectHumidity) > 0.01 {
		t.Fatalf("expected Humidity %.2f, got %.2f", expectHumidity, got)
	}

	if got, _ := d.PM2p5(); abs(got-expectPM2p5) > 0.01 {
		t.Fatalf("expected PM2.5 %.2f, got %.2f", expectPM2p5, got)
	}

	if got, _ := d.CO2(); int(got) != expectCO2 {
		t.Fatalf("expected CO2 %d, got %d", expectCO2, int(got))
	}

	if got, _ := d.VOC(); int(got) != expectVOC {
		t.Fatalf("expected VOC %d, got %d", expectVOC, int(got))
	}

	if got, _ := d.NOX(); int(got) != expectNOX {
		t.Fatalf("expected NOX %d, got %d", expectNOX, int(got))
	}

	if got, _ := d.Luminosity(); abs(got-expectLuminosity) > 0.1 {
		t.Fatalf("expected Luminosity %.2f, got %.2f", expectLuminosity, got)
	}

	if got, _ := d.MeasurementSequenceNumber(); int(got) != expectSequenceNumber {
		t.Fatalf("expected SequenceNumber %d, got %d", expectSequenceNumber, int(got))
	}

	if got, _ := d.CalibrationInProgress(); got != expectFlagCalibrationInProgress {
	 	t.Fatalf("expected FlagCalibrationInProgress %v, got %v", expectFlagCalibrationInProgress, got)
	}

	if got, _ := d.MACAddress(); !bytes.Equal(got, expectMAC) {
		t.Fatalf("expected MAC 0x%X, got 0x%X", expectMAC, got)
	}
}

func TestInvalidData(t *testing.T) {
	// Test vector: https://docs.ruuvi.com/communication/bluetooth-advertisements/data-format-6#case-invalid-values
	// 0x068000FFFFFFFFFFFFFFFFFFFFFFXXFFFFFFFFFF
	data := []byte{
		0x06, 0x80, 0x00, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
	}

	const (
		expectDataFormat = 6
		expectSequenceNumber = 255
		expectFlagCalibrationInProgress = true
	)
	expectMAC := []byte{0xFF, 0xFF, 0xFF}

	d, err := NewDataFormat6(data)
	if err != nil {
		t.Fatalf("unexpected error constructing DataFormat6: %e", err)
	}

	gotDataFormat := d.DataFormat()
	if expectDataFormat != gotDataFormat {
		t.Fatalf("expected DataFormat %d, got %d", expectDataFormat, gotDataFormat)
	}

	if _, err := d.Temperature(); err == nil {
		t.Fatal("expected error from Temperature")
	}
	if _, err := d.Pressure(); err == nil {
		t.Fatal("expected error from Pressure")
	}
	if _, err := d.Humidity(); err == nil {
		t.Fatal("expected error from Humidity")
	}
	if _, err := d.PM2p5(); err == nil {
		t.Fatal("expected error from PM2p5")
	}
	if _, err := d.CO2(); err == nil {
		t.Fatal("expected error from CO2")
	}
	if _, err := d.VOC(); err == nil {
		t.Fatal("expected error from VOC")
	}
	if _, err := d.NOX(); err == nil {
		t.Fatal("expected error from NOX")
	}
	if _, err := d.Luminosity(); err == nil {
		t.Fatal("expected error from Luminosity")
	}
	if got, _ := d.MeasurementSequenceNumber(); got != expectSequenceNumber {
		t.Fatalf("expected MeasurementSequenceNumber %d, got %d", expectSequenceNumber, got)
	}

	if got, _ := d.CalibrationInProgress(); got != expectFlagCalibrationInProgress {
	 	t.Fatalf("expected CalibrationInProgress %v, got %v", expectFlagCalibrationInProgress, got)
	}

	if got, _ := d.MACAddress(); !bytes.Equal(got, expectMAC) {
		t.Fatalf("expected MAC 0x%X, got 0x%X", expectMAC, got)
	}
}

func TestRawData(t *testing.T) {
	data := []byte{
		0x06, 0x17, 0x0C, 0x56,
		0x68, 0xC7, 0x9E, 0x00,
		0x70, 0x00, 0xC9, 0x05,
		0x01, 0xD9, 0x7F, 0xCD,
		0x00, 0x4C, 0x88, 0x4F,
	}
	df6, _ := NewDataFormat6(data)

	if !cmp.Equal(data, df6.RawData()) {
		t.Fatal("RawData return different data than what was put in!")
	}
}
func TestDataModifiedWithoutCopy(t *testing.T) {
	data := []byte{
		0x06, 0x17, 0x0C, 0x56,
		0x68, 0xC7, 0x9E, 0x00,
		0x70, 0x00, 0xC9, 0x05,
		0x01, 0xD9, 0x7F, 0xCD,
		0x00, 0x4C, 0x88, 0x4F,
	}
	df6, _ := NewDataFormat6(data)

	data[2] = 0x00
	data[5] = 0xFF

	b := df6.RawData()
	if b[2] != 0x00 {
		t.Fatal("underlying data not modified without Copy()")
	}
	if b[5] != 0xFF {
		t.Fatal("underlying data not modified without Copy()")
	}
}

func TestDataNotModifiedWithCopy(t *testing.T) {
	data := []byte{
		0x06, 0x17, 0x0C, 0x56,
		0x68, 0xC7, 0x9E, 0x00,
		0x70, 0x00, 0xC9, 0x05,
		0x01, 0xD9, 0x7F, 0xCD,
		0x00, 0x4C, 0x88, 0x4F,
	}
	df6, _ := NewDataFormat6(data)
	df6.Copy()

	data[2] = 0x00
	data[5] = 0xFF

	b := df6.RawData()
	if b[2] != 0x0C {
		t.Fatal("underlying data modified after calling Copy()")
	}
	if b[5] != 0xC7 {
		t.Fatal("underlying data modified after calling Copy()")
	}
}

func TestErrorReturnedOnBadInput(t *testing.T) {
	wrongDataFormat := []byte{
		0x02, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02, 0xCA, 0x0B, 0x53,
	}
	tooShort := []byte{
		0x06, 0x29, 0x1A, 0x1E, 0xCE, 0x1E, 0xFC,
		0x18, 0xF9, 0x42, 0x02,
	}

	if _, err := NewDataFormat6(wrongDataFormat); err == nil {
		t.Fatal("No error from wrong data format")
	}
	if _, err := NewDataFormat6(tooShort); err == nil {
		t.Fatal("No error from too short data")
	}
}
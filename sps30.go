//--------------------------------------------------------------------------------------------------
//
// Copyright (c) 2021 Patricio Díaz
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
// associated documentation files (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial
// portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
// BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
// DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
//--------------------------------------------------------------------------------------------------

package sps30

import i2c "github.com/d2r2/go-i2c"

// AirQualityReading data readed from sensor
type AirQualityReading struct {
	// all big-endian, IEEE754 float value
	MassPM1             float32 // Mass Concentration PM1.0 [μg/m3]
	MassPM25            float32 // Mass Concentration PM2.5 [μg/m3]
	MassPM4             float32 // Mass Concentration PM4.0 [μg/m3]
	MassPM10            float32 // Mass Concentration PM10 [μg/m3]
	NumberPM05          float32 // Number Concentration PM0.5 [#/cm3]
	NumberPM1           float32 // Number Concentration PM1.0 [#/cm3]
	NumberPM25          float32 // Number Concentration PM2.5 [#/cm3]
	NumberPM4           float32 // Number Concentration PM4.0 [#/cm3]
	NumberPM10          float32 // Number Concentration PM10 [#/cm3]
	TypicalParticleSize float32 // Typical Particle Size [μm]
}

// SPS30 represents a sps30 sensor
type SPS30 struct {
	i2c *i2c.I2C
}

// NewSPS30 creates a new sensor object
func NewSPS30(i2c *i2c.I2C) *SPS30 {
	return &SPS30{i2c: i2c}
}

// ReadArticleCode reads sensor article code
func (v *SPS30) ReadArticleCode() (string, error) {
	var code []byte

	b0, err := readFromAddress(v.i2c, []byte{0xD0, 0x25}, 47)
	if err != nil {
		return "", err
	}
	lg.Debugf("Read: %v\n", b0)

	for i := 0; i <= len(b0); i++ {
		if b0[i] == 0 {
			break
		}

		if (i % 3) != 2 {
			code = append(code, b0[i])
		}
	} // for ...

	return string(code), nil
}

// ReadSerial reads sensor serial
func (v *SPS30) ReadSerial() (string, error) {
	var serial []byte
	b0, err := readFromAddress(v.i2c, []byte{0xD0, 0x33}, 47)
	if err != nil {
		return "", err
	}

	for i := 0; i <= len(b0); i++ {
		if (i % 3) != 2 {
			serial = append(serial, b0[i])
		}
	} // for ...
	return string(serial), nil
}

// ReadCleaningInterval reads cleaning interval from sensor
func (v *SPS30) ReadCleaningInterval() (int64, error) {
	b0, err := readFromAddress(v.i2c, []byte{0x80, 0x04}, 6)
	if err != nil {
		return -1, err
	}

	return int64(b0[4]) + (int64(b0[3]) << 8) + (int64(b0[1]) << 16) + (int64(b0[0]) << 24), nil
}

// StartMeasurement starts measurement
func (v *SPS30) StartMeasurement() error {
	err := writeToAddress(v.i2c, []byte{0x00, 0x10, 0x03, 0x00, crc([]byte{0x03, 0x00})})
	return err
}

// StopMeasurement stops measurements
func (v *SPS30) StopMeasurement() error {
	err := writeToAddress(v.i2c, []byte{0x01, 0x04})
	return err
}

// ReadMeasurement reads measurements
func (v *SPS30) ReadMeasurement() (*AirQualityReading, error) {
	b0, err := readFromAddress(v.i2c, []byte{0x03, 0x00}, 60)
	if err != nil {
		return nil, err
	}
	// lg.Debugf("ReadMeasurement: %X", b0)

	var measurement AirQualityReading = AirQualityReading{
		MassPM1:             byteArrayToFloat32([]byte{b0[0], b0[1], b0[3], b0[4]}),
		MassPM25:            byteArrayToFloat32([]byte{b0[6], b0[7], b0[9], b0[10]}),
		MassPM4:             byteArrayToFloat32([]byte{b0[12], b0[13], b0[15], b0[16]}),
		MassPM10:            byteArrayToFloat32([]byte{b0[18], b0[19], b0[21], b0[22]}),
		NumberPM05:          byteArrayToFloat32([]byte{b0[24], b0[25], b0[27], b0[27]}),
		NumberPM1:           byteArrayToFloat32([]byte{b0[30], b0[31], b0[33], b0[34]}),
		NumberPM25:          byteArrayToFloat32([]byte{b0[36], b0[37], b0[39], b0[40]}),
		NumberPM4:           byteArrayToFloat32([]byte{b0[42], b0[43], b0[45], b0[46]}),
		NumberPM10:          byteArrayToFloat32([]byte{b0[48], b0[49], b0[51], b0[52]}),
		TypicalParticleSize: byteArrayToFloat32([]byte{b0[54], b0[55], b0[57], b0[58]}),
	}

	return &measurement, nil
}

// ReadDataReady reads data-ready flag
// 0x00: no new measurements available
// 0x01: new measurements ready to read
func (v *SPS30) ReadDataReady() int {
	b0, err := readFromAddress(v.i2c, []byte{0x02, 0x02}, 3)
	if err != nil {
		return -1
	}

	if b0[1] == 0x01 {
		return 1
	}

	return 0
}

// Reset sends reset command
func (v *SPS30) Reset() error {
	err := writeToAddress(v.i2c, []byte{0xD3, 0x04})
	return err
}

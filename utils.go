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

import (
	"encoding/binary"
	"math"

	i2c "github.com/d2r2/go-i2c"
	"github.com/sigurn/crc8"
)

var CRC8_SPS30 = crc8.Params{0x31, 0xFF, false, false, 0x00, 0x00, "CRC-8/SPS30"}
var table *crc8.Table = crc8.MakeTable(CRC8_SPS30)

func readFromAddress(i2c *i2c.I2C, addr []byte, numBytes int) ([]byte, error) {
	b0 := make([]byte, numBytes)

	_, err := i2c.WriteBytes(addr)
	if err != nil {
		lg.Errorf("Error sending code: %X", addr)
		return nil, err
	}

	_, err = i2c.ReadBytes(b0)
	if err != nil {
		lg.Errorf("Reading response for: %X", addr)
		return nil, err
	}

	return b0, nil
} // readFromAddress ...

func writeToAddress(i2c *i2c.I2C, data []byte) error {
	_, err := i2c.WriteBytes(data)
	if err != nil {
		lg.Errorf("Error sending code: %X", data)
		return err
	}

	return nil
} // writeToAddress

func crc(data []byte) byte {
	return crc8.Checksum(data, table)
} // crc ...

func byteArrayToInt64(data []byte) int64 {
	return int64(data[4]) + (int64(data[3]) << 8) + (int64(data[1]) << 16) + (int64(data[0]) << 24)
} // calcInt64 ...

func byteArrayToFloat32(data []byte) float32 {
	a0 := binary.BigEndian.Uint32(data)
	a1 := math.Float32frombits(a0)
	return a1
} // calcFloat64 ...

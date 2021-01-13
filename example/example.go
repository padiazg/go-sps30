package main

import (
	"fmt"
	"time"

	i2c "github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
	sps30 "github.com/padiazg/go-sps30"
)

var lg = logger.NewPackageLogger("main",
	logger.InfoLevel,
)

func main() {
	defer logger.FinalizeLogger()

	// Create new connection to i2c-bus on 1 line with address 0x69.
	// Use i2cdetect utility to find device address over the i2c-bus
	i2c, err := i2c.NewI2C(0x69, 1)
	if err != nil {
		lg.Error("Creating connection...", err)
	}
	defer i2c.Close()

	logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	logger.ChangePackageLogLevel("sps30", logger.InfoLevel)

	sensor := sps30.NewSPS30(i2c)

	// Read serial
	serial, err := sensor.ReadSerial()
	lg.Infof("Serial: %s", serial)

	// Read cleaning interval
	var cleaningInterval int64
	cleaningInterval, err = sensor.ReadCleaningInterval()
	lg.Infof("interval: %v", cleaningInterval)

	// Start air quality measurement, go out of idle-mode
	err = sensor.StartMeasurement()
	time.Sleep(1 * time.Second)

	// Read if there is new data available
	var dataReady int = sensor.ReadDataReady()
	// lg.Infof("data-ready: %v", dataReady)

	if dataReady == 1 {
		// Read measurements
		m, err := sensor.ReadMeasurement()
		if err != nil {
			lg.Errorf("read-measurement: %v", err)
			return
		}
		// print to console
		formatMeasurementHuman(m)
	}

	// Stop measurement, go to idle-mode again
	err = sensor.StopMeasurement()
}

// Formats data for human reading
func formatMeasurementHuman(m *sps30.AirQualityReading) {
	fmt.Printf("pm0.5 count: %8s\n", fmt.Sprintf("%4.3f", m.NumberPM05))
	fmt.Printf("pm1   count: %8s ug: %6s\n", fmt.Sprintf("%4.3f", m.NumberPM1), fmt.Sprintf("%2.3f", m.MassPM1))
	fmt.Printf("pm2.5 count: %8s ug: %6s\n", fmt.Sprintf("%4.3f", m.NumberPM25), fmt.Sprintf("%2.3f", m.MassPM25))
	fmt.Printf("pm4   count: %8s ug: %6s\n", fmt.Sprintf("%4.3f", m.NumberPM4), fmt.Sprintf("%2.3f", m.MassPM4))
	fmt.Printf("pm10  count: %8s ug: %6s\n", fmt.Sprintf("%4.3f", m.NumberPM10), fmt.Sprintf("%2.3f", m.MassPM10))
	fmt.Printf("pm_typ: %4.3f\n", m.TypicalParticleSize)
}

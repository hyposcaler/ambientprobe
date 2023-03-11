package main

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
)

var (
	ambientTemperature = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ambient_temperature",
		Help: "The current Ambient Temperature",
	})
	ambientPressure = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ambient_pressure",
		Help: "The current Ambient Pressure",
	})
	ambientHumidity = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ambient_humidity",
		Help: "The current Ambient Humidity",
	})
)

func main() {

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		panic(err)
	}

	b := bme280.New(d)
	err = b.Init()

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			t, p, h, err := b.EnvData()
			if err != nil {
				panic(err)
			}
			ambientTemperature.Set(t)
			ambientPressure.Set(p)
			ambientHumidity.Set(h)
			// fmt.Printf("Temp: %fC, Press: %fhPa, Hum: %f%%\n", t, p, h)
			time.Sleep(15 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2323", nil)

}

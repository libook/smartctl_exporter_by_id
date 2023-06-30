package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	devices    []*Device
	results    = sync.Map{}
	collectors = map[string]*prometheus.GaugeVec{}

	_flags = Flags{}

	AppName   = "smartctl_exporter"
	Version   = ""
	BuildDate = ""
)

func WithMetrics(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wg := sync.WaitGroup{}
		wg.Add(len(devices))
		for _, d := range devices {
			go func(d *Device) {
				results.Store(d.Name, GetAll(d))
				wg.Done()
			}(d)
		}
		wg.Wait()

		for _, d := range devices {
			if r, ok := results.Load(d.Name); ok {

				r := r.(*Result)
				passed := 0
				if r.Passed {
					passed = 1
				}
				collectors["device_status"].
					WithLabelValues(d.Name, r.ModelName, r.SerialNumber, r.FirmwareVersion, d.LabelPath).
					Set(float64(passed))

				for k, v := range r.Attributes {
					c, ok := collectors[k]
					if !ok {
						c = promauto.NewGaugeVec(
							prometheus.GaugeOpts{
								Namespace: "smartctl",
								Name:      k,
							},
							[]string{"device", "label_path"},
						)
						collectors[k] = c
					}
					c.WithLabelValues(d.Name, d.LabelPath).Set(v)
				}
			}
		}

		handler.ServeHTTP(w, r)
	}
}

func main() {
	_flags.init()

	if *_flags.Version {
		fmt.Printf(
			"%s\n"+
				"Version: \t%s\n"+
				"Build date: \t%s\n",
			AppName,
			Version,
			BuildDate)
		return
	}

	devices = GetDevices()

	deviceStatusCollector := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "smartctl",
			Name:      "device_status",
			Help:      "Device Status",
		},
		[]string{
			"device",
			"model_name",
			"serial_number",
			"firmware_version",
			"label_path",
		},
	)
	collectors["device_status"] = deviceStatusCollector

	promHandler := promhttp.Handler()

	hf := WithMetrics(promHandler)

	if !*_flags.disableAuth {
		hf = BasicAuth(hf)
	}

	http.HandleFunc(*_flags.Path, hf)

	addr := fmt.Sprintf("%s:%d", *_flags.Address, *_flags.Port)

	log.Printf("Listen: %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

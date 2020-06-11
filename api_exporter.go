package main

import (
	"api_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log2 "github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
)

func main() {
	var (
		port        = kingpin.Flag("port", "The port to be listened.").Short('p').Default(":9080").String()
		metricsPath = kingpin.Flag("metrics-path", "Path under which to expose metrics.").Default("/metrics").String()
		accessPath  = kingpin.Flag("file", "Path to the access file.").Short('f').Default("").String()
	)

	kingpin.Parse()

	exporter, err := collector.NewAccessCollector(*accessPath)

	if err != nil {
		//TODO
	}

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(exporter)

	gatherers := prometheus.Gatherers{
		reg,
	}

	handler := promhttp.HandlerFor(gatherers,
		promhttp.HandlerOpts{
			ErrorLog:      log2.NewErrorLogger(),
			ErrorHandling: promhttp.ContinueOnError,
		})

	http.Handle(*metricsPath, handler)

	log.Fatal(http.ListenAndServe(*port, nil))

}

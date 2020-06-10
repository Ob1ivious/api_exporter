package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	http.Handle(*metricsPath, promhttp.Handler())

	log.Fatal(http.ListenAndServe(*port, nil))

}

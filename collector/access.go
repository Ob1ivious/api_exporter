package collector

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const accessSubsystem = "access"

var (
	accessLabelNames = []string{"api"}

	apiResponseDesc = prometheus.NewDesc(prometheus.BuildFQName(namespace, accessSubsystem, "api_response_time"),
		"The response time of api request.",
		accessLabelNames, nil,
	)
)

type accessCollector struct {
	sync.Mutex

	apiResponseDesc *prometheus.Desc
	accessFilePath  string

	mux *http.ServeMux
}

func NewAccessCollector(accessFilePath string) (prometheus.Collector, error) {
	e := &accessCollector{
		apiResponseDesc: apiResponseDesc,
		accessFilePath:  accessFilePath,
	}

	e.mux = http.NewServeMux()
	return e, nil
}

func (c *accessCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.apiResponseDesc
}

func (c *accessCollector) Collect(ch chan<- prometheus.Metric) {

	for api, time := range c.analysisAccessLog(c.accessFilePath) {
		ch <- prometheus.MustNewConstMetric(
			c.apiResponseDesc,
			prometheus.CounterValue,
			time,
			api,
		)
	}

}

func (c *accessCollector) analysisAccessLog(filePath string) map[string]float64 {

	analysisResult := make(map[string]float64)

	file, err := os.Open(filePath)

	if err != nil {
		//TODO 文件打开异常处理
	}

	defer func() {
		if err = file.Close(); err != nil {
			//TODO 文件关闭异常处理
		}
	}()

	line := bufio.NewScanner(file)

	for line.Scan() {
		slice := strings.Split(line.Text(), " ")
		analysisResult[slice[6]], err = strconv.ParseFloat(slice[len(slice)-1], 64)

		if err != nil {
			//TODO 字符串转float异常处理
		}
	}

	err = line.Err()

	if err != nil {
		//TODO 文件读取行异常处理
	}

	return analysisResult
}

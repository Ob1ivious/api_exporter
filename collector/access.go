package collector

import (
	"api_exporter/utils"
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
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

	for api, consumeTime := range c.analysisAccessLog(c.accessFilePath) {
		ch <- prometheus.MustNewConstMetric(
			c.apiResponseDesc,
			prometheus.GaugeValue,
			consumeTime,
			api,
		)
	}

}

func (c *accessCollector) analysisAccessLog(filePath string) map[string]float64 {

	totalResult := make(map[string][]float64)
	avgResult := make(map[string]float64)

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

		url := slice[6]
		state := slice[8]
		timeStr := slice[len(slice)-1]
		requestTimeStr := slice[3]

		loc, _ := time.LoadLocation("Local")
		requestTime, _ := time.ParseInLocation("2/Jan/2006:15:04:05", requestTimeStr[1:], loc)

		now := time.Now().Unix()

		slice = strings.Split(url, "?")
		api := slice[0]

		consumeTime, err := strconv.ParseFloat(timeStr[1:len(timeStr)-1], 64)


		if state == "200" && api[0] == '/' && (now-requestTime.Unix() < 3600*24) {

			if len(api) > 6 && api[0:6] == "/image" {
				api = "/image"
			}

			totalResult[api] = append(totalResult[api], consumeTime)
		}

		if err != nil {
			//TODO 字符串转float异常处理
		}
	}

	for api, total := range totalResult {
		avgResult[api] = utils.Round(utils.Avg(total), 3)
	}

	err = line.Err()

	if err != nil {
		//TODO 文件读取行异常处理
	}

	return avgResult
}

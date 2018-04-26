package main

import (
	//"client_golang/prometheus/push"
	"fmt"
	"os"
	"time"

	"github.com/prometheus/client_golang/api/prometheus"
	prm "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"golang.org/x/net/context"
)

// ConsumerUtilization helper string to construct a query for consumer utilization
const ConsumerUtilization string = "rabbitmq_queue_consumer_utilisation"

// ReadyMessages helper string to construct a query for messages ready in the queue
const ReadyMessages string = "rabbitmq_queue_messages_ready"

const totalPublished string = "rabbitmq_queue_messages_published_total"
const totalDelivered string = "rabbitmq_queue_messages_delivered_total"

const consumers string = "rabbitmq_queue_consumers"

//upper time bound to wait for retrieving the metrics
const upperTimeBound time.Duration = 500 * time.Millisecond

var (
	numberOfReplicas = prm.NewGauge(prm.GaugeOpts{
		Name: "numberOfReplicas",
		Help: "number of replicas",
	})

	queueIncreaseRate = prm.NewGauge(prm.GaugeOpts{
		Name: "queueIncreaseRate",
		Help: "the rate of increase from the queue",
	})
)

type prometheusMonitor struct {
	client          *prometheus.Client
	prometheusURL   string
	queueToObserve  string
	metricToObserve string
}

//NewPrometheusMonitor creates a prometheusMonitor with a prometheusClient initialized
func NewPrometheusMonitor(prometheusURL string, queueToObserve string, metricToObserve string) (pm *prometheusMonitor, err error) {
	if len(os.Getenv("PROMETHEUS_HOST")) > 0 && len(os.Getenv("PROMETHEUS_PORT")) > 0 {
		prometheusURL = fmt.Sprintf("http://%v:%v", os.Getenv("PROMETHEUS_HOST"), os.Getenv("PROMETHEUS_PORT"))
	}
	prometheusClient, err := createPrometheusClient(prometheusURL)

	if err != nil {
		return nil, err
	}
	monitor := prometheusMonitor{}
	monitor.client = prometheusClient
	monitor.queueToObserve = queueToObserve
	monitor.metricToObserve = metricToObserve
	return &monitor, nil
}

func createPrometheusClient(prometheusURL string) (client *prometheus.Client, err error) {

	promClient, err := prometheus.New(prometheus.Config{Address: prometheusURL})
	client = &promClient
	if err != nil {
		//fmt.Println("%v", err)
		//		panic(err)
		return nil, err
	}
	return client, nil
}

// GetRate returns the rate of messages published for the last 1 minute
func (pm prometheusMonitor) GetRate() (total float64, err error) {
	results, err := pm.queryMetrics(fmt.Sprintf("rate(%v{queue='%v'}[1m])", totalPublished, pm.queueToObserve))
	if err != nil {
		return 0.0, err
	}
	return getNumericalValue(results), err
}

//helper methods
func (pm prometheusMonitor) queryMetrics(query string) (val model.Vector, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), upperTimeBound)
	defer cancel()
	q := prometheus.NewQueryAPI(*pm.client)

	//context call
	value, err := q.Query(ctx, query, time.Now())

	if err != nil {
		return nil, err
	}

	return value.(model.Vector), nil
}

func getNumericalValue(model model.Vector) (numerical float64) {
	var answer float64
	for _, val := range model {
		answer = float64(val.Value)
	}
	return answer
}

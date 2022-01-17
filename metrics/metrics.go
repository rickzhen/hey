package metrics

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rickzhen/hey/snapshot"
	"github.com/rickzhen/hey/utils"
)

var rwlock sync.RWMutex

type Miner struct {
	stopChan chan struct{}
	inerval  time.Duration
	host     string
	port     int
	start    time.Duration
	snapshot *snapshot.Report

	fastest      prometheus.Gauge
	slowest      prometheus.Gauge
	average      prometheus.Gauge
	rps          prometheus.Gauge
	numRes       prometheus.Gauge
	avgConn      prometheus.Gauge
	avgDNS       prometheus.Gauge
	avgReq       prometheus.Gauge
	avgRes       prometheus.Gauge
	avgDelay     prometheus.Gauge
	fastestConn  prometheus.Gauge
	slowestConn  prometheus.Gauge
	fastestDns   prometheus.Gauge
	slowestDns   prometheus.Gauge
	fastestReq   prometheus.Gauge
	slowestReq   prometheus.Gauge
	fastestRes   prometheus.Gauge
	slowestRes   prometheus.Gauge
	fastestDelay prometheus.Gauge
	slowestDelay prometheus.Gauge
}

func NewMiner() *Miner {
	return &Miner{}
}

func (m *Miner) Init() {
	m.stopChan = make(chan struct{})
	m.inerval = 500
	m.host = "localhost"
	m.port = 1010
	m.snapshot = &snapshot.Report{}
	m.rps = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_rps",
		Help: "hey request of per second"},
	)
	m.numRes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_num_res",
		Help: "hey num of all requests"},
	)
	m.average = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_average_latency",
		Help: "hey average lantency of all requests"},
	)
	m.fastest = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest",
		Help: "hey minimum lantency of all requests"},
	)
	m.slowest = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowtest",
		Help: "hey maximum lantency of all requests"},
	)
	m.avgConn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_avg_conn",
		Help: "hey average lantency of connection setup(DNS lookup + Dial up)"},
	)
	m.avgDNS = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_avg_dns",
		Help: "hey average lantency of dns lookup"},
	)
	m.avgReq = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_avg_req",
		Help: "hey average lantency of request \"write\""},
	)
	m.avgRes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_avg_res",
		Help: "hey average lantency of response \"read\""},
	)
	m.avgDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_avg_delay",
		Help: "hey average lantency between response and request"},
	)

	m.fastestConn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest_conn",
		Help: "hey minimum lantency of connection setup(DNS lookup + Dial up)"},
	)
	m.slowestConn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowest_conn",
		Help: "hey maximum lantency of connection setup(DNS lookup + Dial up)"},
	)
	m.fastestDns = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest_dns",
		Help: "hey minimum lantency of dns lookup"},
	)
	m.slowestDns = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowest_dns",
		Help: "hey maximum lantency of dns lookup"},
	)
	m.fastestReq = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest_req",
		Help: "hey minimum lantency of request \"write\""},
	)
	m.slowestReq = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowest_req",
		Help: "hey maximum lantency of request \"write\""},
	)
	m.fastestRes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest_res",
		Help: "hey minimum lantency of response \"read\""},
	)
	m.slowestRes = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowest_res",
		Help: "hey maximum lantency of response \"read\""},
	)
	m.fastestDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_fastest_delay",
		Help: "hey minimum lantency between response and request"},
	)
	m.slowestDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "hey_slowest_delay",
		Help: "hey maximum lantency between response and request"},
	)
	prometheus.MustRegister(m.rps, m.numRes, m.average, m.fastest, m.slowest, m.avgConn, m.avgDNS, m.avgReq, m.avgRes, m.avgDelay)
	prometheus.MustRegister(m.fastestConn, m.slowestConn, m.fastestDns, m.slowestDns, m.fastestReq, m.slowestReq, m.fastestRes, m.slowestRes, m.fastestDelay, m.slowestDelay)

}

func (m *Miner) Run() {
	m.Init()
	m.start = utils.Now()
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", m.host, m.port), nil)
		if err != nil {
			log.Println("error:", err.Error())
			return
		}
	}()
	for {
		rwlock.RLock()
		m.rps.Set(m.snapshot.Rps)
		m.numRes.Set(float64(m.snapshot.NumRes))
		m.average.Set(m.snapshot.Average)
		m.fastest.Set(m.snapshot.Fastest)
		m.slowest.Set(m.snapshot.Slowest)
		m.avgConn.Set(m.snapshot.AvgConn)
		m.avgDNS.Set(m.snapshot.AvgDNS)
		m.avgReq.Set(m.snapshot.AvgReq)
		m.avgRes.Set(m.snapshot.AvgRes)
		m.avgDelay.Set(m.snapshot.AvgDelay)
		m.fastestConn.Set(m.snapshot.ConnMax)
		m.slowestConn.Set(m.snapshot.ConnMin)
		m.fastestDns.Set(m.snapshot.DnsMax)
		m.slowestDns.Set(m.snapshot.DnsMin)
		m.fastestReq.Set(m.snapshot.ReqMax)
		m.slowestReq.Set(m.snapshot.ReqMin)
		m.fastestRes.Set(m.snapshot.ResMax)
		m.slowestRes.Set(m.snapshot.ResMin)
		m.fastestDelay.Set(m.snapshot.DelayMax)
		m.slowestDelay.Set(m.snapshot.DelayMin)
		rwlock.RUnlock()
		time.Sleep(time.Millisecond * m.inerval)
		select {
		case <-m.stopChan:
			return
		default:
			continue
		}
	}
}

func (m *Miner) Stop() {
	m.stopChan <- struct{}{}
}

func (m *Miner) SetSnapshot(snap *snapshot.Report) {
	rwlock.Lock()
	m.snapshot = snap
	rwlock.Unlock()
}

package snapshot

import "time"

type Report struct {
	AvgTotal float64
	Fastest  float64
	Slowest  float64
	Average  float64
	Rps      float64

	AvgConn  float64
	AvgDNS   float64
	AvgReq   float64
	AvgRes   float64
	AvgDelay float64
	ConnMax  float64
	ConnMin  float64
	DnsMax   float64
	DnsMin   float64
	ReqMax   float64
	ReqMin   float64
	ResMax   float64
	ResMin   float64
	DelayMax float64
	DelayMin float64

	Lats        []float64
	ConnLats    []float64
	DnsLats     []float64
	ReqLats     []float64
	ResLats     []float64
	DelayLats   []float64
	Offsets     []float64
	StatusCodes []int

	Total time.Duration

	ErrorDist      map[string]int
	StatusCodeDist map[int]int
	SizeTotal      int64
	SizeReq        int64
	NumRes         int64

	LatencyDistribution []LatencyDistribution
	Histogram           []Bucket
}
type LatencyDistribution struct {
	Percentage int
	Latency    float64
}

type Bucket struct {
	Mark      float64
	Count     int
	Frequency float64
}

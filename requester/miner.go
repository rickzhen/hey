package requester

import (
	"time"

	"github.com/rickzhen/hey/snapshot"
	"github.com/rickzhen/hey/utils"
)

func (w *Work) runMiner() {
	// cap := min(w.N, maxRes)
	snapshot := &snapshot.Report{
		Fastest:     0,
		ConnMax:     0,
		DnsMax:      0,
		ReqMax:      0,
		ResMax:      0,
		DelayMax:    0,
		StatusCodes: make([]int, 0),
	}
	interval := 0
	var avgTotal float64
	var avgConn float64
	var avgDelay float64
	var avgDNS float64
	var avgReq float64
	var avgRes float64
	var succNum int64
	var sizeTotal int64
	var numRes int64
	for res := range w.snapshots {
		if time.Duration((utils.Now()-w.start).Milliseconds())/time.Duration(w.MetricsInterval) != time.Duration(interval) {
			snapshot.Rps = float64(numRes) / (float64(w.MetricsInterval) / 1000)
			w.miner.SetSnapshot(snapshot)
			avgTotal = 0
			avgConn = 0
			avgDelay = 0
			avgDNS = 0
			avgReq = 0
			avgRes = 0
			succNum = 0
			sizeTotal = 0
			numRes = 0
			interval++
		}
		snapshot.NumRes++
		numRes++
		if res.err != nil {
			// snapshot.ErrorDist[res.err.Error()]++
			continue
		} else {
			avgTotal += res.duration.Seconds()
			avgConn += res.connDuration.Seconds()
			avgDelay += res.delayDuration.Seconds()
			avgDNS += res.dnsDuration.Seconds()
			avgReq += res.reqDuration.Seconds()
			avgRes += res.resDuration.Seconds()
			succNum += 1
			switch snapshot.Fastest {
			case 0:
				snapshot.Fastest = res.duration.Seconds()
			default:
				snapshot.Fastest = fmin(snapshot.Fastest, res.duration.Seconds())
			}
			snapshot.Slowest = fmax(snapshot.Slowest, res.duration.Seconds())
			switch snapshot.ConnMax {
			case 0:
				snapshot.ConnMax = res.connDuration.Seconds()
			default:
				snapshot.ConnMax = fmin(snapshot.ConnMax, res.connDuration.Seconds())
			}
			snapshot.ConnMin = fmax(snapshot.ConnMin, res.connDuration.Seconds())
			switch snapshot.DnsMax {
			case 0:
				snapshot.DnsMax = res.dnsDuration.Seconds()
			default:
				snapshot.DnsMax = fmin(snapshot.DnsMax, res.dnsDuration.Seconds())
			}
			snapshot.DnsMin = fmax(snapshot.DnsMin, res.dnsDuration.Seconds())
			switch snapshot.ReqMax {
			case 0:
				snapshot.ReqMax = res.reqDuration.Seconds()
			default:
				snapshot.ReqMax = fmin(snapshot.ReqMax, res.reqDuration.Seconds())
			}
			snapshot.ReqMin = fmax(snapshot.ReqMin, res.reqDuration.Seconds())
			switch snapshot.ResMax {
			case 0:
				snapshot.ResMax = res.resDuration.Seconds()
			default:
				snapshot.ResMax = fmin(snapshot.ResMax, res.resDuration.Seconds())
			}
			snapshot.ResMin = fmax(snapshot.ResMin, res.resDuration.Seconds())
			switch snapshot.DelayMax {
			case 0:
				snapshot.DelayMax = res.delayDuration.Seconds()
			default:
				snapshot.DelayMax = fmin(snapshot.DelayMax, res.delayDuration.Seconds())
			}
			snapshot.DelayMin = fmax(snapshot.DelayMin, res.delayDuration.Seconds())
		}
		if res.contentLength > 0 {
			snapshot.SizeTotal += res.contentLength
			sizeTotal += res.contentLength
		}
		snapshot.Total = now() - w.start
		snapshot.Average = avgTotal / float64(succNum)
		snapshot.AvgConn = avgConn / float64(succNum)
		snapshot.AvgDelay = avgDelay / float64(succNum)
		snapshot.AvgDNS = avgDNS / float64(succNum)
		snapshot.AvgReq = avgReq / float64(succNum)
		snapshot.AvgRes = avgRes / float64(succNum)
		snapshot.SizeReq = sizeTotal / succNum
	}
}

func fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
func fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

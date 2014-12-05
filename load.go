package gorilla

import "fmt"
import linuxproc "github.com/c9s/goprocinfo/linux"

type LoadAverage struct {
	last1m, last5m, last15m float64
        critical, warning float64
}

func (l *LoadAverage) Usage() float64 {
	loadAverage, _ := linuxproc.ReadLoadAvg("/proc/loadavg")
	cpuInfo, _ := linuxproc.ReadCPUInfo("/proc/cpuinfo")

	l.last1m = loadAverage.Last1Min / float64(cpuInfo.NumCore())
	l.last5m = loadAverage.Last5Min / float64(cpuInfo.NumCore())
	l.last15m = loadAverage.Last15Min / float64(cpuInfo.NumCore())

	return l.last1m
}

func (l *LoadAverage) Ranking() string {
	return fmt.Sprintf("1-minute load average/core is %f", l.last1m)
}

func (l *LoadAverage) Collect(queue chan *Metric) {

	metric := NewMetric()

	metric.Service = "load"

        usage := l.Usage()
	metric.Value = usage
	metric.Description = l.Ranking()

        switch {
                case usage > l.critical:
                        metric.State = "critical"
                case usage > l.warning:
                        metric.State = "warning"
        }

        queue <- metric
}

func NewLoadAverage(warning, critical float64) *LoadAverage {
	return &LoadAverage{warning: warning, critical: critical}
}

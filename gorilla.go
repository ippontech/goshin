package gorilla

import (
	"fmt"
	"github.com/bigdatadev/goryman"
	"time"
)

type Metric struct {
	Service, Description, State string
	Value                interface{}
}

func NewMetric() *Metric{
        return &Metric{State: "ok"}
}

type Threshold struct {
        Warning, Critical float64
}

func NewThreshold() *Threshold {
        return &Threshold {}
}

type Gorilla struct {
	Address      string
	EventHost    string
	Interval     int
	Tag          []string
	Ttl          float32
	Ifaces       map[string]bool
	IgnoreIfaces map[string]bool
        Thresholds   map[string]*Threshold
}

func NewGorilla() *Gorilla {
	return &Gorilla{
                Thresholds: make(map[string]*Threshold),
        }
}

func (g *Gorilla) Start() {
	fmt.Print("Gare aux goriiillllleeeees!\n\n\n")

	cputime := NewCPUTime()
	memoryusage := NewMemoryUsage()
	loadaverage := NewLoadAverage()
	netstats := NewNetStats(g.Ifaces, g.IgnoreIfaces)

	fmt.Printf("Gorilla will report each %d seconds\n", g.Interval)

        // channel size has to be large enough
        // to allow Gorilla send all metrics to Riemann
        // in g.Interval
        var collectQueue chan *Metric = make(chan *Metric, 100)

        go g.Report(collectQueue)

	ticker := time.NewTicker(time.Second * time.Duration(g.Interval))

	for t := range ticker.C {
		fmt.Println("Tick at ", t)
                go cputime.Collect(collectQueue)
		go memoryusage.Collect(collectQueue)
		go loadaverage.Collect(collectQueue)
		go netstats.Collect(collectQueue)

                go g.Report(collectQueue)
	}
}

func (g *Gorilla) EnforceState(metric *Metric) {

        threshold, present := g.Thresholds[metric.Service]

        if present {
                value := metric.Value

                // TODO threshold checking
                // only for int and float type
                switch {
                        case value.(float64) > threshold.Critical:
                                metric.State = "critical"
                        case value.(float64)> threshold.Warning:
                                metric.State = "warning"
                        default:
                                metric.State = "ok"
                }
        }
}


func (g *Gorilla) Report(reportQueue chan *Metric) {

        c := goryman.NewGorymanClient(g.Address)
        err := c.Connect()

        if err != nil {
                fmt.Println("Can not connect to host")
        } else {

                more := true

                for more {
                        select {
                        case metric := <- reportQueue:
                                g.EnforceState(metric)
                                err := c.SendEvent(&goryman.Event{
                                        Metric:      metric.Value,
                                        Ttl:         g.Ttl,
                                        Service:     metric.Service,
                                        Description: metric.Description,
                                        Tags:        g.Tag,
                                        Host:        g.EventHost,
                                        State:       metric.State})

                                if err != nil {
                                        fmt.Println("something does wrong:", err)
                                }
                        default:
                                more = false
                        }
                }
        }

        defer c.Close()
}

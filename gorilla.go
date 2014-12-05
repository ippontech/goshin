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

type Gorilla struct {
	Address      string
	CheckCPU     bool
	EventHost    string
	Interval     int
	Tag          []string
	Ttl          float32
	Ifaces       map[string]bool
	IgnoreIfaces map[string]bool
        CpuWarning, CpuCritical float64
        LoadWarning, LoadCritical float64
}

func NewGorilla() *Gorilla {
	return &Gorilla{}
}

func (g *Gorilla) Start() {
	fmt.Print("Gare aux goriiillllleeeees!\n\n\n")

	cputime := NewCPUTime(g.CpuWarning, g.CpuCritical)
	memoryusage := NewMemoryUsage()
	loadaverage := NewLoadAverage(g.LoadWarning, g.LoadCritical)
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


func (g *Gorilla) Report(collectQueue chan *Metric) {

        c := goryman.NewGorymanClient(g.Address)
        err := c.Connect()

        if err != nil {
                fmt.Println("Can not connect to host")
        } else {

                more := true

                for more {
                        select {
                        case metric := <- collectQueue:
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

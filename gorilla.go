package gorilla

import (
	"fmt"
	"github.com/bigdatadev/goryman"
	"time"
)

type Metric struct {
	service, description string
	value                interface{}
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
}

func NewGorilla() *Gorilla {
	return &Gorilla{}
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
                                        Metric:      metric.value,
                                        Ttl:         g.Ttl,
                                        Service:     metric.service,
                                        Description: metric.description,
                                        Tags:        g.Tag,
                                        Host:        g.EventHost,
                                        State:       "ok"})

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

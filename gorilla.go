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

func Instance() *Gorilla {
	return &Gorilla{}
}

func (g *Gorilla) Start() {
	fmt.Print("Gare aux goriiillllleeeees!\n\n\n")

	cputime := NewCPUTime()
	memoryusage := NewMemoryUsage()
	loadaverage := NewLoadAverage()
	netstats := NewNetStats(g.Ifaces, g.IgnoreIfaces)

	reporter := func(metric *Metric) {

		c := goryman.NewGorymanClient(g.Address)
		err := c.Connect()

		if err != nil {
			fmt.Println("can not connect to host")
		} else {
			c.SendEvent(&goryman.Event{
				Metric:      metric.value,
				Ttl:         g.Ttl,
				Service:     metric.service,
				Description: metric.description,
				Tags:        g.Tag,
				Host:        g.EventHost,
				State:       "ok"})
		}

		defer c.Close()
	}

	fmt.Printf("Gorilla will report each %d seconds\n", g.Interval)
	ticker := time.NewTicker(time.Second * time.Duration(g.Interval))

	for t := range ticker.C {
		fmt.Println("Tick at ", t)
		go cputime.Report(reporter)
		go memoryusage.Report(reporter)
		go loadaverage.Report(reporter)
		go netstats.Report(reporter)
	}
}

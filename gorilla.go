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
	Address   string
	CheckCPU  bool
	EventHost string
	Interval  int
	Tag       string
	Ttl       float64
}

func Instance() *Gorilla {
	return &Gorilla{}
}

func (g *Gorilla) Start() {
	fmt.Print("Gare aux goriiillllleeeees!\n")

	cputime := new(CPUTime)
	memoryusage := new(MemoryUsage)
	loadaverage := new(LoadAverage)
	netstats := NewNetStats()

	reporter := func(metric *Metric) {

		c := goryman.NewGorymanClient(g.Address)
		err := c.Connect()

		if err != nil {
			fmt.Println("can not connect to host")
		} else {
			c.SendEvent(&goryman.Event{
				Metric:      metric.value,
				Ttl:         10,
				Service:     metric.service,
				Description: metric.description,
				State:       "ok"})
		}

		defer c.Close()
	}

	ticker := time.NewTicker(time.Second * 2)

	for t := range ticker.C {
		fmt.Println("Tick at ", t)
		go cputime.Report(reporter)
		go memoryusage.Report(reporter)
		go loadaverage.Report(reporter)
		go netstats.Report(reporter)
	}
}

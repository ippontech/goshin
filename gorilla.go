package gorilla

import (
        "fmt"
        "time"
        "github.com/bigdatadev/goryman"
    )


type Metric struct {
        service, description string
        value interface {}
}

type Gorilla struct {
        Address string
        CheckCPU bool
        EventHost string
        Interval int
        Tag string
        Ttl float64
}

func Instance() *Gorilla {
        return &Gorilla{}
}


func (g *Gorilla) Start() {
        fmt.Print("Gare aux goriiillllleeeees!\n")


        cputime := new (CPUTime)
        memoryusage := new (MemoryUsage)
        loadaverage := new (LoadAverage)
        netstats := NewNetStats()

 //       var metricQueue chan *Metric = make(chan *Metric, 1)

        c := goryman.NewGorymanClient(g.Address)
        c.Connect()
        reporter := func (metric *Metric) {


                c.SendEvent(&goryman.Event{
                        Metric: metric.value,
                        Ttl: 10,
                        Service: metric.service,
                        Description: metric.description,
                        State: "ok"})

//                defer c.Close()
        }


        ticker := time.NewTicker(time.Second * 2)

        for  t:= range ticker.C {
                fmt.Println("Tick at ", t)
                go cputime.Report(reporter)
                go memoryusage.Report(reporter)
                go loadaverage.Report(reporter)
                go netstats.Report(reporter)
        }
}


func (g *Gorilla) Report(metricQueue chan *Metric) {

        c := goryman.NewGorymanClient(g.Address)
        err := c.Connect()
        if err != nil {
                panic(fmt.Sprintf("Can not open connection to Riemann %s. Check your configuration.", g.Address))
        }


        for {
                metric := <- metricQueue

                go send(c, metric)
        }
}

func send(c *goryman.GorymanClient, metric *Metric) {

        //fmt.Println("Send metric : ", metric)
        err := c.SendEvent(&goryman.Event{
                Metric: metric.value,
                Ttl: 10,
                Service: metric.service,
                Description: metric.description,
                State: "ok"})


        if err != nil {
                fmt.Println("wtf?")
        }
}

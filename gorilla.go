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

        reporter := func (metric *Metric) {

                c := goryman.NewGorymanClient(g.Address)
                err := c.Connect()
                if err == nil {
                        err := c.SendEvent(&goryman.Event{
                                Metric: metric.value,
                                Ttl: float32(g.Ttl),
                                Host: g.EventHost,
                                Service: metric.service,
                                Description: metric.description,
                                State: "ok"})

                        if err != nil {
                                fmt.Println("wtf?")
                        }

                        defer c.Close()
                } else {
                        panic(fmt.Sprintf("Can not open connection to Riemann %s. Check your configuration.", g.Address))
                }
        }

        cputime := new (CPUTime)
        memoryusage := new (MemoryUsage)
        loadaverage := new (LoadAverage)

        for i := 0 ; i < 1000 ; i++ {
                cputime.Report(reporter)
                memoryusage.Report(reporter)
                loadaverage.Report(reporter)
                time.Sleep(time.Duration(g.Interval) * time.Second)
        }
}

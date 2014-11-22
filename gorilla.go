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

        var metricQueue chan *Metric = make(chan *Metric, 10)


        go g.Report(metricQueue)


        ticker := time.NewTicker(time.Second * 5)

        for t:= range ticker.C {
                fmt.Println("Tick at ", t)
                go cputime.Report(metricQueue)
                go memoryusage.Report(metricQueue)
                go loadaverage.Report(metricQueue)
        }
}


func (g *Gorilla) Report(metricQueue chan *Metric) {

        buffer := make([]*Metric, 3)

        for {

                for index, _ := range buffer {
                        buffer[index] = <- metricQueue
                }

                c := goryman.NewGorymanClient(g.Address)
                err := c.Connect()
                if err == nil {

                        fmt.Println("Send Riemann events")

                        for _, metric := range buffer {
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

                        }

                        defer c.Close()
                } else {
                        panic(fmt.Sprintf("Can not open connection to Riemann %s. Check your configuration.", g.Address))
                }
        }
}

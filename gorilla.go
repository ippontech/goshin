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
}

func Instance() *Gorilla {
        return &Gorilla{}
}


func (g *Gorilla) Start() {
        fmt.Print("Gare aux gorillllles!\n")

        reporter := func (metric *Metric) {

                c := goryman.NewGorymanClient(g.Address)
                c.Connect()

                c.SendEvent(&goryman.Event{
                        Metric: metric.value,
                        Ttl: 10,
                        Service: metric.service,
                        Description: metric.description,
                        State: "ok"})

                defer c.Close()
        }

        cputime := new (CPUTime)
        memoryusage := new (MemoryUsage)
        loadaverage := new (LoadAverage)

        for i := 0 ; i < 1000 ; i++ {
                cputime.Report(reporter)
                memoryusage.Report(reporter)
                loadaverage.Report(reporter)
                time.Sleep(5 * time.Second)
        }
}

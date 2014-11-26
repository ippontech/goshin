package gorilla

import (
        "fmt"
    )
import linuxproc "github.com/c9s/goprocinfo/linux"



type NetStats struct {
        last, actual  map[string]linuxproc.NetworkStat
}

func (n *NetStats) Store() {

        // no copy for map on "=" with GO
        for k, v := range n.actual {
                n.last[k] = v
        }

        netStat,_ := linuxproc.ReadNetworkStat("/proc/net/dev")


        for _, ifaceStat := range netStat {

                ifaceName := ifaceStat.Iface

                // store new value
                n.actual[ifaceName] = ifaceStat
        }
}

func buildMetric(iface string, name string, actual uint64, last uint64) *Metric {
        m := new(Metric)
        m.service = fmt.Sprintf("%s %s", iface, name)

        diff := actual - last
        m.value = float64(diff) / float64(5)

        return m
}


func (n  NetStats) Report(metricQueue chan *Metric) {
        n.Store()

        fmt.Println("start report")

        for ifaceName, actualStat := range n.actual {

                lastStat := n.last[ifaceName]

                metricQueue <- buildMetric(ifaceName, "rx bytes", actualStat.RxBytes, lastStat.RxBytes)
                metricQueue <- buildMetric(ifaceName, "rx packets", actualStat.RxPackets, lastStat.RxPackets)
                metricQueue <- buildMetric(ifaceName, "rx errs", actualStat.RxErrs, lastStat.RxErrs)
                metricQueue <- buildMetric(ifaceName, "rx drop", actualStat.RxDrop, lastStat.RxDrop)
                metricQueue <- buildMetric(ifaceName, "rx frame", actualStat.RxFrame, lastStat.RxFrame)
                metricQueue <- buildMetric(ifaceName, "rx compressed", actualStat.RxCompressed, lastStat.RxCompressed)
                metricQueue <- buildMetric(ifaceName, "rx muticast", actualStat.RxMulticast, lastStat.RxMulticast)


                metricQueue <- buildMetric(ifaceName, "tx bytes", actualStat.TxBytes, lastStat.TxBytes)
                metricQueue <- buildMetric(ifaceName, "tx packets", actualStat.TxPackets, lastStat.TxPackets)
                metricQueue <- buildMetric(ifaceName, "tx errs", actualStat.TxErrs, lastStat.TxErrs)
                metricQueue <- buildMetric(ifaceName, "tx drop", actualStat.TxDrop, lastStat.TxDrop)
                metricQueue <- buildMetric(ifaceName, "tx fifo", actualStat.TxFifo, lastStat.TxFifo)
                metricQueue <- buildMetric(ifaceName, "tx colls", actualStat.TxColls, lastStat.TxColls)
                metricQueue <- buildMetric(ifaceName, "tx carrier", actualStat.TxCarrier, lastStat.TxCarrier)
                metricQueue <- buildMetric(ifaceName, "tx compressed", actualStat.TxCompressed, lastStat.TxCompressed)


        }

        fmt.Println("stop report")
}

// Act as constructor
func NewNetStats() *NetStats {
        return &NetStats{ 
                        last:make(map[string]linuxproc.NetworkStat),
                        actual: make(map[string]linuxproc.NetworkStat)}
}

package gorilla

import (
	"fmt"
	"time"
)
import linuxproc "github.com/c9s/goprocinfo/linux"

type NetStats struct {
	last, actual         map[string]linuxproc.NetworkStat
	lastTime, actualTime time.Time
	ignoreIfaces, ifaces map[string]bool
}

func (n *NetStats) Store() {

	n.lastTime = n.actualTime

	// no copy for map on "=" with GO
	for k, v := range n.actual {
		n.last[k] = v
	}

	netStat, _ := linuxproc.ReadNetworkStat("/proc/net/dev")

	for _, ifaceStat := range netStat {

		ifaceName := ifaceStat.Iface

		// store new value
		n.actual[ifaceName] = ifaceStat
	}

	n.actualTime = time.Now()

}

func buildMetric(iface string, name string, actual uint64, last uint64, interval float64) *Metric {
	m := new(Metric)
	m.service = fmt.Sprintf("%s %s", iface, name)

	diff := actual - last
	m.value = float64(diff) / interval

	return m
}

func (n *NetStats) candidateIfaces() []string {

	keys := make([]string, 0, len(n.actual))

	for k, _ := range n.actual {
		_, include := n.ifaces[k]
		_, exclude := n.ignoreIfaces[k]

		if len(n.ifaces) != 0 {
			if include && !exclude {
				keys = append(keys, k)
			}
		} else {
			if !exclude {
				keys = append(keys, k)
			}
		}

	}

	return keys
}

func (n *NetStats) Report(f func(*Metric)) {
	n.Store()

	// first run or
	// no interface
	if len(n.last) == 0 {
		return
	}

	interval := float64(n.actualTime.Sub(n.lastTime).Seconds())

	for _, ifaceName := range n.candidateIfaces() {

		lastStat := n.last[ifaceName]
		actualStat := n.actual[ifaceName]

		f(buildMetric(ifaceName, "rx bytes", actualStat.RxBytes, lastStat.RxBytes, interval))
		f(buildMetric(ifaceName, "rx packets", actualStat.RxPackets, lastStat.RxPackets, interval))
		f(buildMetric(ifaceName, "rx errs", actualStat.RxErrs, lastStat.RxErrs, interval))
		f(buildMetric(ifaceName, "rx drop", actualStat.RxDrop, lastStat.RxDrop, interval))
		f(buildMetric(ifaceName, "rx frame", actualStat.RxFrame, lastStat.RxFrame, interval))
		f(buildMetric(ifaceName, "rx compressed", actualStat.RxCompressed, lastStat.RxCompressed, interval))
		f(buildMetric(ifaceName, "rx muticast", actualStat.RxMulticast, lastStat.RxMulticast, interval))

		f(buildMetric(ifaceName, "tx bytes", actualStat.TxBytes, lastStat.TxBytes, interval))
		f(buildMetric(ifaceName, "tx packets", actualStat.TxPackets, lastStat.TxPackets, interval))
		f(buildMetric(ifaceName, "tx errs", actualStat.TxErrs, lastStat.TxErrs, interval))
		f(buildMetric(ifaceName, "tx drop", actualStat.TxDrop, lastStat.TxDrop, interval))
		f(buildMetric(ifaceName, "tx fifo", actualStat.TxFifo, lastStat.TxFifo, interval))
		f(buildMetric(ifaceName, "tx colls", actualStat.TxColls, lastStat.TxColls, interval))
		f(buildMetric(ifaceName, "tx carrier", actualStat.TxCarrier, lastStat.TxCarrier, interval))
		f(buildMetric(ifaceName, "tx compressed", actualStat.TxCompressed, lastStat.TxCompressed, interval))

	}
}

// Act as constructor
func NewNetStats(ifaces, ignoreIfaces map[string]bool) *NetStats {
	return &NetStats{
		last:         make(map[string]linuxproc.NetworkStat),
		actual:       make(map[string]linuxproc.NetworkStat),
		ignoreIfaces: ignoreIfaces,
		ifaces:       ifaces,
	}
}

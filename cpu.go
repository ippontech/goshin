package goshin

import (
	"fmt"
	"github.com/tjgq/broadcast"
	"os/exec"
)
import linuxproc "github.com/c9s/goprocinfo/linux"

type CPUTime struct {
	last, actual linuxproc.CPUStat
}

func (c *CPUTime) Store() {
	c.last = c.actual
	stat, _ := linuxproc.ReadStat("/proc/stat")
	c.actual = stat.CPUStatAll
}

func (c *CPUTime) Used() uint64 {
	var used = (c.actual.User + c.actual.Nice + c.actual.System) - (c.last.User + c.last.Nice + c.last.System)

	return used
}

func (c *CPUTime) IOWait() uint64 {
	return c.actual.IOWait - c.last.IOWait
}

func (c *CPUTime) Total() uint64 {
	// should IOWait be considered as idle?
	return c.Used() + (c.actual.Idle + c.actual.IOWait) - (c.last.Idle + c.actual.IOWait)
}

func (c *CPUTime) Usage() float64 {
	var fraction float64 = float64(c.Used()) / float64(c.Total())
	return fraction
}

func (c *CPUTime) IOWaitUsage() float64 {
	var fraction float64 = float64(c.IOWait()) / float64(c.Total())
	return fraction
}

func (c *CPUTime) Ranking() string {
	out, _ := exec.Command("sh", "-c", "ps -eo pcpu,pid,comm | sort -nrb -k1 | head -10").Output()

	s := string(out[:])

	return fmt.Sprint("user+nice+system\n\n", s)
}

func (c *CPUTime) Collect(queue chan *Metric, listener *broadcast.Listener) {

	for {
		<-listener.Ch

		c.Store()

		if c.last.User == 0 {
			// nothing stored yet
			// so no metric to send
			continue
		}
		cpu := NewMetric()

		cpu.Service = "cpu"
		cpu.Value = c.Usage()
		cpu.Description = c.Ranking()

		queue <- cpu

		cpuwait := NewMetric()

		cpuwait.Service = "cpuwait"
		cpuwait.Value = c.IOWaitUsage()
		cpuwait.Description = c.Ranking()

		queue <- cpuwait
	}
}

func NewCPUTime() *CPUTime {
	return &CPUTime{}
}

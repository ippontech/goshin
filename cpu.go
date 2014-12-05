package gorilla

import (
	"fmt"
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

func (c *CPUTime) Total() uint64 {
	return c.Used() + c.actual.Idle - c.last.Idle
}

func (c *CPUTime) Usage() float64 {
	var fraction float64 = float64(c.Used()) / float64(c.Total())
	return fraction
}

func (c *CPUTime) Ranking() string {
	out, _ := exec.Command("sh", "-c", "ps -eo pcpu,pid,comm | sort -nrb -k1 | head -10").Output()

	s := string(out[:])

	return fmt.Sprint("user+nice+system\n\n", s)
}

func (c *CPUTime) Collect(queue chan *Metric) {

	c.Store()

        if c.last.User ==  0 {
                // nothing stored yet
                // so no metric to send
                return
        }
	metric := new(Metric)

	metric.Service = "cpu"
	metric.Value = c.Usage()
	metric.Description = c.Ranking()

	queue <- metric
}

func NewCPUTime() *CPUTime {
	return &CPUTime{}
}

package goshin

import "fmt"
import "os/exec"
import linuxproc "github.com/c9s/goprocinfo/linux"

type MemoryUsage struct {
	total, free, buffers, cache, swap uint64
}

func (m *MemoryUsage) Usage() float64 {
	memInfo, _ := linuxproc.ReadMemInfo("/proc/meminfo")

	m.free = memInfo.MemFree + memInfo.Buffers + memInfo.Cached
	m.total = memInfo.MemTotal

	return float64(1 - float64(m.free)/float64(m.total))
}

func (m *MemoryUsage) Ranking() string {
	out, _ := exec.Command("sh", "-c", "ps -eo pmem,pid,comm | sort -nrb -k1 | head -10").Output()

	s := string(out[:])

	return fmt.Sprint("used\n\n", s)
}

func (m *MemoryUsage) Collect(queue chan *Metric) {

	metric := NewMetric()

	metric.Service = "memory"
	metric.Value = m.Usage()
	metric.Description = m.Ranking()

	queue <- metric
}

func NewMemoryUsage() *MemoryUsage {
	return &MemoryUsage{}
}

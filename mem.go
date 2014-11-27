package gorilla

import "fmt"
import "os/exec"
import linuxproc "github.com/c9s/goprocinfo/linux"

type MemoryUsage struct {
        total, free, buffers, cache, swap uint64
}

func (m* MemoryUsage) Usage() float64 {
        memInfo, _ := linuxproc.ReadMemInfo("/proc/meminfo")

        m.free = memInfo["MemFree"] + memInfo["Buffers"] + memInfo["Cached"]
        m.total = memInfo["MemTotal"]

        return float64(1 - float64(m.free) / float64(m.total))
}

func (m* MemoryUsage) Ranking() string {
        out, _ :=  exec.Command("sh", "-c", "ps -eo pmem,pid,comm | sort -nrb -k1 | head -10").Output()

        s := string(out[:])

       return fmt.Sprint("used\n\n", s)
}

func (m* MemoryUsage) Report(f func(*Metric))  {

        metric := new(Metric)

        metric.service = "memory"
        metric.value = m.Usage()
        metric.description = m.Ranking()

        f(metric)
}

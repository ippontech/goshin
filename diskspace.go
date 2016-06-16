package goshin

import "fmt"
import "os/exec"
import "strings"
import "strconv"
import "github.com/tjgq/broadcast"

type DiskSpace struct {
}

func (m *DiskSpace) Collect(queue chan *Metric, listener *broadcast.Listener) {
	for {
		<-listener.Ch

		out, _ := exec.Command("sh", "-c", "df -P").Output()

		lines := strings.Split(string(out), "\n")

		for _, line := range lines[1:] {
			fields := strings.Fields(line)

			if len(fields) == 0 {
				continue
			}

			if !strings.Contains(fields[0], "/") {
				continue
			}

			var capacity float64
			fields[4] = strings.Replace(fields[4], "%", "", -1)
			capacity, _ = strconv.ParseFloat(fields[4], 64)

			metric := NewMetric()
			metric.Service = fmt.Sprint("disk ", fields[5])
			metric.Value = capacity / 100
			metric.Description = fmt.Sprint(capacity, "% used")

			queue <- metric
		}
	}
}

func NewDiskSpace() *DiskSpace {
	return &DiskSpace{}
}

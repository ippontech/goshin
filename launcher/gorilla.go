package main

import "github.com/pariviere/gorilla"
import "flag"
import "fmt"
import "os"
import "github.com/vharitonsky/iniflags"
import "strings"

var (
	hostname, _ = os.Hostname()

	hostPtr         = flag.String("host", "localhost", "Riemann host")
	portPtr         = flag.Int("port", 5555, "Riemann port")
	eventHostPtr    = flag.String("event-host", hostname, "Event hostname")
	intervalPtr     = flag.Int("interval", 5, "Seconds between updates")
	tagPtr          = flag.String("tag", "", "Tag to add to events")
	ttlPtr          = flag.Float64("ttl", 10, "TTL for events")
	ifacesPtr       = flag.String("interfaces", "", "Interfaces to monitor")
	ignoreIfacesPtr = flag.String("ignore-interfaces", "lo", "Interfaces to ignore (default: lo)")
        cpuWarningPtr   = flag.Float64("cpu-warning", 0.9, "CPU warning threshold (fraction of total jiffies")
        cpuCriticalPtr  = flag.Float64("cpu-critical", 0.95, "CPU critical threshold (fraction of total jiffies")
        loadWarningPtr   = flag.Float64("load-warning", 3, "Load warning threshold (load average / core")
        loadCriticalPtr  = flag.Float64("load-critical", 8, "Load critical threshold (load average / core)")
)

func main() {

	iniflags.Parse()

	app := gorilla.NewGorilla()

	app.Address = fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
	app.EventHost = *eventHostPtr
	app.Interval = *intervalPtr

	if len(*tagPtr) != 0 {
		app.Tag = strings.Split(*tagPtr, ",")
	}

	app.Ttl = float32(*ttlPtr)

	ifaces := make(map[string]bool)

	if len(*ifacesPtr) != 0 {
		for _, iface := range strings.Split(*ifacesPtr, ",") {
			ifaces[iface] = true
		}
	}
	app.Ifaces = ifaces

	ignoreIfaces := make(map[string]bool)

	if len(*ignoreIfacesPtr) != 0 {
		for _, ignoreIface := range strings.Split(*ignoreIfacesPtr, ",") {
			ignoreIfaces[ignoreIface] = true
		}
	}
	app.IgnoreIfaces = ignoreIfaces


        cpuThreshold := gorilla.NewThreshold()
        cpuThreshold.Critical = *cpuCriticalPtr
        cpuThreshold.Warning = *cpuWarningPtr

        app.Thresholds["cpu"] = cpuThreshold


        loadThreshold := gorilla.NewThreshold()
        loadThreshold.Critical = *loadCriticalPtr
        loadThreshold.Warning = *loadWarningPtr

        app.Thresholds["load"] = loadThreshold

	app.Start()
}

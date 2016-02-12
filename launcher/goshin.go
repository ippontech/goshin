package main

import "github.com/ippontech/goshin"
import "flag"
import "fmt"
import "os"
import "github.com/vharitonsky/iniflags"
import "strings"

var (
	hostname, _ = os.Hostname()

	hostPtr           = flag.String("host", "localhost", "Riemann host")
	portPtr           = flag.Int("port", 5555, "Riemann port")
	eventHostPtr      = flag.String("event-host", hostname, "Event hostname")
	intervalPtr       = flag.Int("interval", 5, "Seconds between updates")
	tagPtr            = flag.String("tag", "", "Tag to add to events")
	ttlPtr            = flag.Float64("ttl", 10, "TTL for events")
	ifacesPtr         = flag.String("interfaces", "", "Interfaces to monitor")
	ignoreIfacesPtr   = flag.String("ignore-interfaces", "lo", "Interfaces to ignore (default: lo)")
	devicesPtr        = flag.String("devices", "", "Devices to monitor")
	ignoreDevicesPtr  = flag.String("ignore-devices", "", "Devices to ignore (default: nil)")
	cpuWarningPtr     = flag.Float64("cpu-warning", 0.9, "CPU warning threshold (fraction of total jiffies)")
	cpuCriticalPtr    = flag.Float64("cpu-critical", 0.95, "CPU critical threshold (fraction of total jiffies)")
	diskWarningPtr    = flag.Float64("disk-warning", 0.9, "Disk warning threshold (fraction of space used)")
	diskCriticalPtr   = flag.Float64("disk-critical", 0.95, "Disk critical threshold (fraction of space used)")
	loadWarningPtr    = flag.Float64("load-warning", 3, "Load warning threshold (load average / core)")
	loadCriticalPtr   = flag.Float64("load-critical", 8, "Load critical threshold (load average / core)")
	memoryWarningPtr  = flag.Float64("memory-warning", 0.85, "Memory warning threshold (fraction of RAM)")
	memoryCriticalPtr = flag.Float64("memory-critical", 0.95, "Memory critical threshold (fraction of RAM)")
	checksPtr         = flag.String("checks", "cpu,load,memory,net,disk", "A list of checks to run")
)

func main() {

	iniflags.Parse()

	app := goshin.NewGoshin()

	app.Address = fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
	app.EventHost = *eventHostPtr
	app.Interval = *intervalPtr

	if len(*tagPtr) != 0 {
		app.Tag = strings.Split(*tagPtr, ",")
	}

	app.Ttl = float32(*ttlPtr)

	// iface
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

	// devices
	devices := make(map[string]bool)

	if len(*devicesPtr) != 0 {
		for _, device := range strings.Split(*devicesPtr, ",") {
			devices[device] = true
		}
	}
	app.Devices = devices

	ignoreDevices := make(map[string]bool)

	if len(*ignoreDevicesPtr) != 0 {
		for _, ignoreDevice := range strings.Split(*ignoreDevicesPtr, ",") {
			ignoreDevices[ignoreDevice] = true
		}
	}
	app.IgnoreDevices = ignoreDevices

	// threshold
	cpuThreshold := goshin.NewThreshold()
	cpuThreshold.Critical = *cpuCriticalPtr
	cpuThreshold.Warning = *cpuWarningPtr

	app.Thresholds["cpu"] = cpuThreshold

	loadThreshold := goshin.NewThreshold()
	loadThreshold.Critical = *loadCriticalPtr
	loadThreshold.Warning = *loadWarningPtr

	app.Thresholds["load"] = loadThreshold

	memoryThreshold := goshin.NewThreshold()
	memoryThreshold.Critical = *memoryCriticalPtr
	memoryThreshold.Warning = *memoryWarningPtr

	app.Thresholds["memory"] = memoryThreshold

	diskThreshold := goshin.NewThreshold()
	diskThreshold.Critical = *diskCriticalPtr
	diskThreshold.Warning = *diskWarningPtr

	app.Thresholds["disk"] = diskThreshold

	checks := make(map[string]bool)

	if len(*checksPtr) != 0 {
		for _, check := range strings.Split(*checksPtr, ",") {
			checks[check] = true
		}
	}

	app.Checks = checks

	app.Start()
}

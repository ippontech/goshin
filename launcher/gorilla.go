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
)

func main() {

	iniflags.Parse()

	gorilla := gorilla.Instance()

	gorilla.Address = fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
	gorilla.EventHost = *eventHostPtr
	gorilla.Interval = *intervalPtr

	if len(*tagPtr) != 0 {
		gorilla.Tag = strings.Split(*tagPtr, ",")
	}

	gorilla.Ttl = float32(*ttlPtr)

	ifaces := make(map[string]bool)

	if len(*ifacesPtr) != 0 {
		for _, iface := range strings.Split(*ifacesPtr, ",") {
			ifaces[iface] = true
		}
	}
	gorilla.Ifaces = ifaces

	ignoreIfaces := make(map[string]bool)

	if len(*ignoreIfacesPtr) != 0 {
		for _, ignoreIface := range strings.Split(*ignoreIfacesPtr, ",") {
			ignoreIfaces[ignoreIface] = true
		}
	}
	gorilla.IgnoreIfaces = ignoreIfaces

	gorilla.Start()
}

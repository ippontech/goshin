package main

import "github.com/pariviere/gorilla"
import "flag"
import "fmt"
import "os"

func main() {

        hostname,_ := os.Hostname()

        hostPtr := flag.String("host", "localhost", "Riemann host")
        portPtr := flag.Int("port", 5555, "Riemann port")
        eventHostPtr := flag.String("event_host",  hostname, "Event hostname")
        intervalPtr := flag.Int("interval", 5, "Seconds between updates")
        tagPtr := flag.String("tag", "", "Tag to add to events")
        ttlPtr := flag.Float64("ttl", 10, "TTL for events")

        flag.Parse()

        gorilla := gorilla.Instance()

        gorilla.Address = fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
        gorilla.EventHost = *eventHostPtr
        gorilla.Interval = *intervalPtr
        gorilla.Tag = *tagPtr
        gorilla.Ttl = *ttlPtr

        gorilla.Start()
}

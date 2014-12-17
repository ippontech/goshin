# Gorilla!

> I love Riemann but I hate Riemann-tools. You too? Try Gorilla!


Gorilla! aims to be a drop-in alternative to Riemann-tools.

## Why?

By itself, Riemann-tools does the job. The main criticism is the installation procedure 

```
gem install riemann-tools
```

A pretty neat one-line installer... Cool, right? But let's take a look at the dependencies :

 - __faraday__: HTTP/REST API client library
 - __fog__: The Ruby cloud services library
 - __munin-client__: Munin Node client
 - __Nokogiri__: an HTML, XML, SAX, and Reader parser
 - __trollop__: a commandline option parser for Ruby that just gets out of your way
 - __yajl-ruby__: Ruby C bindings to the excellent Yajl JSON stream-based parser library.
 - __beefcake__: A sane Google Protocol Buffers library for Ruby


What's wrong? Actually nothing if you just want to push metrics from your personnal laptop. But things get worse if you have to deploy riemann-tools in a heterogeneous environment :

 - beefcake requires Ruby 1.9.3 and on some (not so) old systems no package is available without adding external packages repositories
 - yajl-ruby and Nokogiri will compile a C extension and for this will require gcc
 

## Gorilla!, a Go alternative to riemann-tools

That's our proposal : use Golang to implement a drop-in alternative to riemann-tools. Just drop one big (sic) binary on your box and launch Gorilla! No external dependencies.

Gorilla! tries to mimic riemann-tools metrics gathering and for now covers the following riemann-tools metric : 

 - riemann-health : cpu, memory and load
 - riemann-net : network usage

Gorilla! also borrows riemann-health and riemann-net flags.

## Usage

To start gathering metrics you just have to launch Gorilla! binary. By default, it will report cpu, memory, load and network metrics to a local Riemann instance.

```
$ ./gorilla
Gare aux goriiillllleeeees!
Gorilla will report each 5 seconds

```
It tries to support the same flags as riemann-health and riemann-net.

```
$ ./gorilla --host=riemann.example.com
```

To check the available flags, check the online help

```
$ ./gorilla --help
Usage of ./gorilla:
  -checks="cpu,load,memory,net,disk": A list of checks to run
  -config="": Path to ini config for using in go flags. May be relative to the current executable path.
  -configUpdateInterval=0: Update interval for re-reading config file set via -config flag. Zero disables config file re-reading.
  -cpu-critical=0.95: CPU critical threshold (fraction of total jiffies
  -cpu-warning=0.9: CPU warning threshold (fraction of total jiffies
  -dumpflags=false: Dumps values for all flags defined in the app into stdout in ini-compatible syntax and terminates the app.
  -event-host="t510i": Event hostname
  -host="localhost": Riemann host
  -ignore-interfaces="lo": Interfaces to ignore (default: lo)
  -interfaces="": Interfaces to monitor
  -interval=5: Seconds between updates
  -load-critical=8: Load critical threshold (load average / core)
  -load-warning=3: Load warning threshold (load average / core
  -memory-critical=0.95: Memory critical threshold (fraction of RAM)
  -memory-warning=0.85: Memory warning threshold (fraction of RAM)
  -port=5555: Riemann port
  -tag="": Tag to add to events
  -ttl=10: TTL for events

```
You can also pass flags via an ini-style configuration file and the `--config` flag

```
$ cat gorilla.ini
host=riemann.example.com
interfaces=wlan0,tun0,eth0
tag=gorilla
interval=1
$ ./gorilla --config gorilla.ini
```

License
-------

Copyright 2012-2014 [Ippon Technologies](http://www.ippon.fr) and [Ippon Hosting](http://www.ippon-hosting.com/)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this application except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.





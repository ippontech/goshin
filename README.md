# Gorilla!

> I love Riemann but I hate Riemann-tools. You too? Try Gorilla!


Gorilla! aims to be a drop in alternative to Riemann-tools.

## Why?

By itself, Riemann-tools does the job. The main criticism is the installation procedure 

```
gem install riemann-tools
```

A pretty neat one-line installer... Cool, right? But let's take a look to dependencies :

 - faraday: HTTP/REST API client library
 - fog: The Ruby cloud services library
 - munin-client: Munin Node client
 - Nokogiri: an HTML, XML, SAX, and Reader parser
 - trollop : a commandline option parser for Ruby that just gets out of your way
 - yajl-ruby : Ruby C bindings to the excellent Yajl JSON stream-based parser library.
 - beefcake : A sane Google Protocol Buffers library for Ruby


What's wrong? Actually nothing if you just want to push metrics from your personnal laptop. But things get worse if you have to deploy riemann-tools in a heterogeneous environment :

 - beefcake requiert Ruby 1.9.3 and on some (not so) old systems no package is available without adding external packages repositories
 - yajl-ruby and Nokogiri will build compile a C extension and for this will require gcc
 

## Gorilla!, a Go alternative to riemann-tools

That's our proposal : use Golang to implement a drop in alternative to riemann-tools. Just drop one big (sic) binary on your box and launch Gorilla! No external dependencies.

Gorilla! try to mimic riemann-tools metrics gathering and for now covers the following riemann-tools metric : 

 - riemann-health : cpu, memory and load
 - riemann-net : network usage

Gorilla! also borrows riemann-health and riemann-net flags.
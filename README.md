This crappy project creates a simple webserver to announce WC availability.
It reads incoming data from a serial port when the WC is in-use and responds accordingly.

It creates a poop.log with start and end time in UNIX Timestamp for each session
Also see `/last` for time since last use for better air quality information.

# Install

* [go](https://golang.org/dl/)
* [serial](https://github.com/tarm/serial) `go get github.com/tarm/serial`

# Enviroment Variables

`DEVICE` defaults to `/dev/ttyUSB0`

`PORT` defaults to `8080`

# API

`/cani/` availability of with time in seconds since last use (if code is `200`)

`/cani/logs` returns the poop.log

`/cani/*` where * is _any_ string

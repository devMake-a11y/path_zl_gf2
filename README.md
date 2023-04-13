go-gelf - GELF Library and Writer for Go [changed]
========================================

This implementation currently supports UDP and TCP as a transport
protocol. TLS is unsupported.

The library provides an API that applications can use to log messages
directly to a Graylog server and an `io.Writer` that can be used to
redirect the standard library's log messages (`os.Stdout`) to a
Graylog server.

[GELF]: http://docs.graylog.org/en/2.2/pages/gelf.html
[syslog]: https://tools.ietf.org/html/rfc5424
[chunking]: http://docs.graylog.org/en/2.2/pages/gelf.html#chunked-gelf


Installing
----------

go-gelf is go get-able:

    go get https://github.com/devMake-a11y/path_zl_gf2

Usage
-----

The easiest way to integrate graylog logging into your go app is by
having your `main` function (or even `init`) call `log.SetOutput()`.
By using an `io.MultiWriter`, we can log to both stdout and graylog -
giving us both centralized and local logs.  (Redundancy is nice).

```golang
package main

import (
	"chiTest/internal/go-zerolog-gelf-2/gelf"
	"database/sql"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"time"
)

func main() {
  var graylogAddr string

  flag.StringVar(&graylogAddr, "graylog", "", "graylog server addr")
  flag.Parse()

	hook, err := gelf.NewUDPWriter(graylogAddr, gelf.ZeroLogParser, gelf.Caller)
	if err != nil {
		panic(err)
	}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Output(hook)

	//Set global logger with graylog hook

	for range time.Tick(time.Millisecond * 200) {
		for range time.Tick(time.Millisecond * 200) {
			log.Info().Str("event", "test").Msg("New event")
			log.Warn().Time("date_event", time.Now()).Str("event", "test").Msg("Event changed")
			log.Err(sql.ErrNoRows).Time("date_event", time.Now()).Str("event", "test").Msg("Event error")
			log.Error().Stack().Err(outer()).Msg("")
		}
	}

  // ...
}
```

License
-------

go-gelf is offered under the MIT license, see LICENSE for details.

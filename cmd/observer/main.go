// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run gen-dotenv.go

package main

import (
	"context"
	"github.com/go-pogo/errors"
	"github.com/roeldev/youless-observer/app/observer"
	"github.com/rs/zerolog"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var unmarshalEnv func(conf *observerapp.Config) error
var loggerOut func() io.Writer

func main() {
	var conf observerapp.Config
	errors.FatalOnErr(unmarshalEnv(&conf))

	log := zerolog.New(loggerOut()).Level(conf.Level)
	if conf.WithTimestamp {
		log = log.With().Timestamp().Logger()
	}

	// collecting metrics is always enabled
	conf.Telemetry.Meter.Enabled = true

	if err := conf.Validate(); err != nil {
		fatalErr(log, err, "invalid configuration", 1)
	}

	app, err := observerapp.New(conf, log)
	if err != nil {
		fatalErr(log, err, "unable to create observer app", 2)
	}

	runCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := app.Run(runCtx); err != nil {
			log.Err(err).Msg("error during run")
		}
	}()
	<-runCtx.Done()

	closeCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := app.Shutdown(closeCtx); err != nil {
			fatalErr(log, err, "error during shutdown", 4)
		}
	}()
	<-closeCtx.Done()
}

func fatalErr(log zerolog.Logger, err error, msg string, code int) {
	log.Err(err).Msg(msg)
	os.Exit(errors.GetExitCodeOr(err, code))
}

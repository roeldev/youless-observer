// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

var unmarshalEnv func(v any) error
var loggerOut func() io.Writer

func main() {
	var conf observerapp.Config
	errors.FatalOnErr(unmarshalEnv(&conf))

	log := zerolog.New(loggerOut()).Level(conf.Level)
	if conf.WithTimestamp {
		log = log.With().Timestamp().Logger()
	}

	app, err := observerapp.New(conf, log)
	if err != nil {
		log.Err(err).Msg("unable to create observer app")
		os.Exit(errors.GetExitCodeOr(err, 2))
	}

	run(log, app.Run, app.Shutdown)
}

func run(log zerolog.Logger, run, shutdown func(context.Context) error) {
	runCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := run(runCtx); err != nil {
			log.Err(err).Msg("error during run")
		}
	}()
	<-runCtx.Done()

	closeCtx, stopFn := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer stopFn()
		if err := shutdown(closeCtx); err != nil {
			log.Err(err).Msg("error during graceful shutdown")
			os.Exit(errors.GetExitCodeOr(err, 4))
		}
	}()
	<-closeCtx.Done()
}

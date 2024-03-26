// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package observerapp

import (
	"context"
	"github.com/go-pogo/buildinfo"
	"github.com/go-pogo/errors"
	youlessclient "github.com/roeldev/youless-client"
	"github.com/roeldev/youless-logger"
	"github.com/roeldev/youless-logger/server"
	"github.com/roeldev/youless-observer"
	"github.com/rs/zerolog"
)

const (
	ErrServerCreateFailure   errors.Msg = "failed to create server"
	ErrClientCreateFailure   errors.Msg = "failed to create client"
	ErrObserverCreateFailure errors.Msg = "failed to create observer"
)

// App is the youless-observer application which runs an Observer that observes
// the YouLess device using a youless.Client. It also runs a server to expose
// health status, build info and metrics endpoints.
type App struct {
	server   *server.Server
	client   *youlessclient.Client
	observer *youlessobserver.Observer
}

func New(conf Config, log zerolog.Logger) (*App, error) {
	var err error
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	var app App
	app.server, err = server.New("observer", conf.Server, log, nil,
		server.WithBuildInfo(buildinfo.New(youlessobserver.Version).
			WithExtra("client_version", youlessclient.Version).
			WithExtra("logger_version", youlesslogger.Version),
		),
		server.WithTelemetryAndPrometheus(conf.Telemetry, conf.Prometheus),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrServerCreateFailure)
	}

	app.client, err = youlessclient.NewClient(conf.YouLess,
		youlessclient.WithLogger(youlessclient.NewLogger(log)),
		youlessclient.WithTracerProvider(app.server.TracerProvider()),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrClientCreateFailure)
	}

	app.observer, err = youlessobserver.NewObserver(app.server.MeterProvider(),
		youlessobserver.WithLogger(&logger{log}),
		youlessobserver.WithMeterReading(conf.Observer.MeterReadingRegisterer, app.client),
		youlessobserver.WithPhaseReading(conf.Observer.PhaseReadingRegisterer, app.client),
	)
	if err != nil {
		return nil, errors.Wrap(err, ErrObserverCreateFailure)
	}

	return &app, nil
}

// Run the app by starting the internal observer and server.
func (app *App) Run(ctx context.Context) error {
	if err := app.observer.Start(); err != nil {
		return err
	}
	return app.server.Run(ctx)
}

// Shutdown the app by stopping the internal observer and server.
func (app *App) Shutdown(ctx context.Context) error {
	return errors.Append(app.observer.Stop(), app.server.Shutdown(ctx))
}

var _ youlessobserver.Logger = (*logger)(nil)

type logger struct{ zl zerolog.Logger }

func (l *logger) Register(name string) {
	l.zl.Debug().
		Str("name", name).
		Msg("observer register")
}

func (l *logger) ObserverStart() { l.zl.Info().Msg("observer starting") }

func (l *logger) ObserverStop() { l.zl.Info().Msg("observer stopped") }

// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/go-pogo/errors"
	"github.com/go-pogo/healthcheck"
	"go.opentelemetry.io/otel/metric"
)

const (
	ErrAlreadyStarted errors.Msg = "observer is already started"
	ErrUnstartedStop  errors.Msg = "cannot stop observer that is not started"
)

type Registerer interface {
	CanRegister() bool
	Register(meter metric.Meter) (Registration, error)
}

type Registration interface {
	LastCheck() time.Time
	Unregister() error
}

var _ healthcheck.HealthCheckerRegisterer = (*Observer)(nil)

type Observer struct {
	log  Logger
	prov metric.MeterProvider

	registerers map[string]*registerer
	started     atomic.Bool
}

type registerer struct {
	Registerer
	registration Registration
}

const panicNilMeterProvider = "youlessobserver: metric.MeterProvider must not be nil"

func NewObserver(prov metric.MeterProvider, opts ...Option) (*Observer, error) {
	if prov == nil {
		panic(panicNilMeterProvider)
	}

	o := Observer{
		prov:        prov,
		registerers: make(map[string]*registerer, len(opts)),
	}
	for _, opt := range opts {
		if err := opt(&o); err != nil {
			return nil, err
		}
	}
	if o.log == nil {
		o.log = NopLogger()
	}
	return &o, nil
}

// Start the Observer by registering all the provided Registerer(s) to a
// metric.Meter created by the provided metric.Provider.
func (o *Observer) Start() error {
	if o.started.Load() {
		return errors.New(ErrAlreadyStarted)
	}

	o.log.ObserverStart()
	o.started.Store(true)
	opts := metric.WithInstrumentationVersion(youlessclient.Version)

	for name, r := range o.registerers {
		reg, err := r.Register(o.prov.Meter(name, opts))
		if err != nil {
			return err
		}
		if reg != nil {
			o.log.Register(name)
			r.registration = reg
		}
	}
	return nil
}

// Stop the Observer by unregistering all the previously registered callbacks.
func (o *Observer) Stop() error {
	if !o.started.Load() {
		return errors.New(ErrUnstartedStop)
	}

	var err error
	for _, r := range o.registerers {
		if r.registration == nil {
			continue
		}

		err = errors.Append(err, r.registration.Unregister())
	}

	o.log.ObserverStop()
	return err
}

func (o *Observer) RegisterHealthCheckers(registerer healthcheck.Registerer) {
	for name, r := range o.registerers {
		if r.registration == nil {
			continue
		}

		registerer.Register(name, r.healthChecker(15*time.Second))
	}
}

func (r *registerer) healthChecker(t time.Duration) healthcheck.HealthCheckerFunc {
	return func(_ context.Context) healthcheck.Status {
		if r.registration == nil {
			return healthcheck.StatusUnknown
		}
		if time.Since(r.registration.LastCheck()) > t {
			return healthcheck.StatusUnhealthy
		}
		return healthcheck.StatusHealthy
	}
}

func must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

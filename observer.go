// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

import (
	"github.com/go-pogo/errors"
	"github.com/roeldev/youless-client"
	"go.opentelemetry.io/otel/metric"
	"sync/atomic"
)

const (
	ErrAlreadyStarted errors.Msg = "observer is already started"
	ErrUnstartedStop  errors.Msg = "cannot stop observer that is not started"
)

type Registerer interface {
	Register(meter metric.Meter) (Registration, error)
}

type Registration interface {
	Unregister() error
}

type Observer struct {
	log  Logger
	prov metric.MeterProvider

	registerers   map[string]Registerer
	registrations []Registration
	started       atomic.Bool
}

const panicNilMeterProvider = "youlessobserver: metric.MeterProvider must not be nil"

func NewObserver(prov metric.MeterProvider, opts ...Option) (*Observer, error) {
	if prov == nil {
		panic(panicNilMeterProvider)
	}

	o := Observer{
		prov:          prov,
		registerers:   make(map[string]Registerer, len(opts)),
		registrations: make([]Registration, 0, len(opts)),
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
	opts := metric.WithInstrumentationVersion(youless.Version)

	for name, reg := range o.registerers {
		r, err := reg.Register(o.prov.Meter(name, opts))
		if err != nil {
			return err
		}
		if r != nil {
			o.log.Register(name)
			o.registrations = append(o.registrations, r)
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
	for _, reg := range o.registrations {
		err = errors.Append(err, reg.Unregister())
	}

	o.log.ObserverStop()
	return err
}

func must[T any](res T, err error) T {
	if err != nil {
		panic(err)
	}
	return res
}

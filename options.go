// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

import (
	youlessclient "github.com/roeldev/youless-client"
)

type Option func(o *Observer) error

func WithLogger(l Logger) Option {
	return func(o *Observer) error {
		o.log = l
		return nil
	}
}

func WithRegisterer(name string, reg Registerer) Option {
	return func(o *Observer) error {
		if reg.CanRegister() {
			o.registerers[name] = &registerer{Registerer: reg}
		}
		return nil
	}
}

func WithMeterReading(reg MeterReadingRegisterer, api youlessclient.API) Option {
	return func(o *Observer) error {
		if api != nil {
			reg.WithAPIClient(api)
		}
		return WithRegisterer(MeterReadingObserverName, &reg)(o)
	}
}

func WithPhaseReading(reg PhaseReadingRegisterer, api youlessclient.API) Option {
	return func(o *Observer) error {
		if api != nil {
			reg.WithAPIClient(api)
		}
		return WithRegisterer(PhaseReadingObserverName, &reg)(o)
	}
}

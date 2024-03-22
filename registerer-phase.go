// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

import (
	"context"
	"github.com/go-pogo/errors"
	youlessclient "github.com/roeldev/youless-client"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const PhaseReadingObserverName = "youless.observer.phase"

var _ Registerer = (*PhaseReadingRegisterer)(nil)

type PhaseReadingRegisterer struct {
	client *youlessclient.Client

	ExcludePower bool
	SinglePhase  bool
}

func NewPhaseReadingRegisterer(client *youlessclient.Client) *PhaseReadingRegisterer {
	var reg PhaseReadingRegisterer
	return reg.WithClient(client)
}

func (reg *PhaseReadingRegisterer) WithClient(client *youlessclient.Client) *PhaseReadingRegisterer {
	reg.client = client
	return reg
}

// Register registers metrics gauges to the provided meter and starts observing
// them by getting phase readings from the client.
func (reg *PhaseReadingRegisterer) Register(meter metric.Meter) (Registration, error) {
	if reg.ExcludePower {
		return nil, nil
	}

	return newPhaseReadingRegistration(reg, meter)
}

var _ Registration = (*phaseReadingRegistration)(nil)

type phaseReadingRegistration struct {
	metric.Registration
	conf PhaseReadingRegisterer

	power1 metric.Int64ObservableGauge
	power2 metric.Int64ObservableGauge
	power3 metric.Int64ObservableGauge

	current1 metric.Float64ObservableGauge
	current2 metric.Float64ObservableGauge
	current3 metric.Float64ObservableGauge

	voltage1 metric.Float64ObservableGauge
	voltage2 metric.Float64ObservableGauge
	voltage3 metric.Float64ObservableGauge
}

func newPhaseReadingRegistration(conf *PhaseReadingRegisterer, meter metric.Meter) (*phaseReadingRegistration, error) {
	reg := phaseReadingRegistration{
		conf: *conf,

		power1: must(meter.Int64ObservableGauge("power1",
			metric.WithDescription("The current imported electricity power in Watt on phase 1."),
			metric.WithUnit("W"),
		)),
		current1: must(meter.Float64ObservableGauge("current1",
			metric.WithDescription("The current imported electricity current in Ampère on phase 1."),
			metric.WithUnit("A"),
		)),
		voltage1: must(meter.Float64ObservableGauge("voltage1",
			metric.WithDescription("The current voltage on phase 1."),
			metric.WithUnit("V"),
		)),
	}

	instruments := make([]metric.Observable, 0, 9)
	instruments = append(instruments, reg.power1, reg.current1, reg.voltage1)

	if !conf.SinglePhase {
		reg.power2 = must(meter.Int64ObservableGauge("power2",
			metric.WithDescription("The current imported electricity power in Watt on phase 2."),
			metric.WithUnit("W"),
		))
		reg.power3 = must(meter.Int64ObservableGauge("power3",
			metric.WithDescription("The current imported electricity power in Watt on phase 3."),
			metric.WithUnit("W"),
		))

		reg.current2 = must(meter.Float64ObservableGauge("current2",
			metric.WithDescription("The current imported electricity current in Ampère on phase 2."),
			metric.WithUnit("A"),
		))
		reg.current3 = must(meter.Float64ObservableGauge("current3",
			metric.WithDescription("The current imported electricity current in Ampère on phase 3."),
			metric.WithUnit("A"),
		))

		reg.voltage2 = must(meter.Float64ObservableGauge("voltage2",
			metric.WithDescription("The current voltage on phase 2."),
			metric.WithUnit("V"),
		))
		reg.voltage3 = must(meter.Float64ObservableGauge("voltage3",
			metric.WithDescription("The current voltage on phase 3."),
			metric.WithUnit("V"),
		))

		instruments = append(
			instruments,
			reg.power2,
			reg.power3,
			reg.current2,
			reg.current3,
			reg.voltage2,
			reg.voltage3,
		)
	}

	var err error
	reg.Registration, err = meter.RegisterCallback(reg.callback, instruments...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &reg, nil
}

func (reg *phaseReadingRegistration) callback(ctx context.Context, observer metric.Observer) error {
	d, err := reg.conf.client.GetPhaseReading(ctx)
	if err != nil {
		// todo: send err to channel
		return err
	}

	attr := attribute.NewSet(
		attribute.Int64("tariff", int64(d.Tariff)),
		attribute.Bool("three-phase", !reg.conf.SinglePhase),
	)

	observer.ObserveInt64(reg.power1, d.Power1, metric.WithAttributeSet(attr))
	observer.ObserveFloat64(reg.current1, d.Current1, metric.WithAttributeSet(attr))
	observer.ObserveFloat64(reg.voltage1, d.Voltage1, metric.WithAttributeSet(attr))

	if !reg.conf.SinglePhase {
		observer.ObserveInt64(reg.power2, d.Power2, metric.WithAttributeSet(attr))
		observer.ObserveInt64(reg.power3, d.Power3, metric.WithAttributeSet(attr))

		observer.ObserveFloat64(reg.current2, d.Current2, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.current3, d.Current3, metric.WithAttributeSet(attr))

		observer.ObserveFloat64(reg.voltage2, d.Voltage2, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.voltage3, d.Voltage3, metric.WithAttributeSet(attr))
	}

	return nil
}

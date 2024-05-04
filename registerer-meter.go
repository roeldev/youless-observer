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
	"sync/atomic"
	"time"
)

const MeterReadingObserverName = "youless.observer.meter"

var _ Registerer = (*MeterReadingRegisterer)(nil)

type MeterReadingRegisterer struct {
	client *youlessclient.Client

	ExcludePower bool
	ExcludeS0    bool
	ExcludeGas   bool
	ExcludeWater bool
}

func NewMeterReadingRegisterer(client *youlessclient.Client) *MeterReadingRegisterer {
	var reg MeterReadingRegisterer
	return reg.WithClient(client)
}

func (reg *MeterReadingRegisterer) WithClient(client *youlessclient.Client) *MeterReadingRegisterer {
	reg.client = client
	return reg
}

// Register registers metrics gauges to the provided meter and starts observing
// them by getting meter readings from the client.
func (reg *MeterReadingRegisterer) Register(meter metric.Meter) (Registration, error) {
	if reg.ExcludePower && reg.ExcludeS0 && reg.ExcludeGas && reg.ExcludeWater {
		return nil, nil
	}

	return newMeterReadingRegistration(reg, meter)
}

var _ Registration = (*meterReadingRegistration)(nil)

type meterReadingRegistration struct {
	metric.Registration
	conf MeterReadingRegisterer
	last atomic.Value

	electricityImport1 metric.Float64ObservableGauge
	electricityImport2 metric.Float64ObservableGauge
	electricityExport1 metric.Float64ObservableGauge
	electricityExport2 metric.Float64ObservableGauge
	netElectricity     metric.Float64ObservableGauge
	power              metric.Int64ObservableGauge

	s0Total metric.Float64ObservableGauge
	s0      metric.Int64ObservableGauge
	gas     metric.Float64ObservableGauge
	water   metric.Float64ObservableGauge
}

func newMeterReadingRegistration(conf *MeterReadingRegisterer, meter metric.Meter) (*meterReadingRegistration, error) {
	instruments := make([]metric.Observable, 0, 10)

	reg := meterReadingRegistration{conf: *conf}
	reg.last.Store(time.Time{})

	if !conf.ExcludePower {
		reg.electricityImport1 = must(meter.Float64ObservableGauge("electricity_import_1",
			metric.WithDescription("The meter reading of total imported low tariff electricity in kWh."),
			metric.WithUnit("kW/h"),
		))
		reg.electricityImport2 = must(meter.Float64ObservableGauge("electricity_import_2",
			metric.WithDescription("The meter reading of total imported high tariff electricity in kWh."),
			metric.WithUnit("kW/h"),
		))
		reg.electricityExport1 = must(meter.Float64ObservableGauge("electricity_export_1",
			metric.WithDescription("The meter reading of total exported low tariff electricity in kWh."),
			metric.WithUnit("kW/h"),
		))
		reg.electricityExport2 = must(meter.Float64ObservableGauge("electricity_export_2",
			metric.WithDescription("The meter reading of total exported high tariff electricity in kWh."),
			metric.WithUnit("kW/h"),
		))
		reg.netElectricity = must(meter.Float64ObservableGauge("net_electricity",
			metric.WithDescription("The total measured electricity which equals (electricity_import_1 + electricity_import_2 - electricity_export_1 - electricity_export_2)."),
			metric.WithUnit("kW/h"),
		))
		reg.power = must(meter.Int64ObservableGauge("power",
			metric.WithDescription("The current imported electricity power in Watt."),
			metric.WithUnit("W"),
		))

		instruments = append(
			instruments,
			reg.electricityImport1,
			reg.electricityImport2,
			reg.electricityExport1,
			reg.electricityExport2,
			reg.netElectricity,
			reg.power,
		)
	}
	if !conf.ExcludeS0 {
		reg.s0Total = must(meter.Float64ObservableGauge("s0_total",
			metric.WithDescription("The total power in kWh measured by the S0 meter."),
			metric.WithUnit("kW/h"),
		))
		reg.s0 = must(meter.Int64ObservableGauge("s0",
			metric.WithDescription("The current electricity power measured in Watt from the S0 meter."),
			metric.WithUnit("W"),
		))
		instruments = append(instruments, reg.s0Total, reg.s0)
	}
	if !conf.ExcludeGas {
		reg.gas = must(meter.Float64ObservableGauge("gas",
			metric.WithDescription("The meter reading of delivered gas (in m3) to client."),
			metric.WithUnit("m3"),
		))
		instruments = append(instruments, reg.gas)
	}
	if !conf.ExcludeWater {
		reg.water = must(meter.Float64ObservableGauge("water",
			metric.WithDescription("The meter reading of delivered water (in m3) to client."),
			metric.WithUnit("m3"),
		))
		instruments = append(instruments, reg.water)
	}

	var err error
	reg.Registration, err = meter.RegisterCallback(reg.callback, instruments...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &reg, nil
}

func (reg *meterReadingRegistration) LastCheck() time.Time {
	return reg.last.Load().(time.Time)
}

func (reg *meterReadingRegistration) callback(ctx context.Context, observer metric.Observer) error {
	d, err := reg.conf.client.GetMeterReading(ctx)
	if err != nil {
		// todo: send err to channel
		return err
	}

	if !reg.conf.ExcludePower && d.Timestamp != 0 {
		attr := attribute.NewSet(attribute.Int64("time_unix_nano", d.Time().UnixNano()))
		observer.ObserveFloat64(reg.electricityImport1, d.ElectricityImport1, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.electricityImport2, d.ElectricityImport2, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.electricityExport1, d.ElectricityExport1, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.electricityExport2, d.ElectricityExport2, metric.WithAttributeSet(attr))
		observer.ObserveFloat64(reg.netElectricity, d.NetElectricity, metric.WithAttributeSet(attr))
		observer.ObserveInt64(reg.power, d.Power, metric.WithAttributeSet(attr))
	}
	if !reg.conf.ExcludeS0 && d.S0Timestamp != 0 {
		attr := attribute.NewSet(attribute.Stringer("timestamp", d.S0Time()))
		observer.ObserveFloat64(reg.s0Total, d.S0Total, metric.WithAttributeSet(attr))
		observer.ObserveInt64(reg.s0, d.S0, metric.WithAttributeSet(attr))
	}
	if !reg.conf.ExcludeGas && d.GasTimestamp != 0 {
		attr := attribute.NewSet(attribute.Stringer("timestamp", d.GasTime()))
		observer.ObserveFloat64(reg.gas, d.Gas, metric.WithAttributeSet(attr))
	}
	if !reg.conf.ExcludeWater && d.WaterTimestamp != 0 {
		attr := attribute.NewSet(attribute.Stringer("timestamp", d.WaterTime()))
		observer.ObserveFloat64(reg.water, d.Water, metric.WithAttributeSet(attr))
	}

	reg.last.Store(time.Now())
	return nil
}

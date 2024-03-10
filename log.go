// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

type Logger interface {
	Register(name string)
	ObserverStart()
	ObserverStop()
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (nopLogger) Register(_ string) {}

func (nopLogger) ObserverStart() {}

func (nopLogger) ObserverStop() {}

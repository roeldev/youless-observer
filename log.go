// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package youlessobserver

type Logger interface {
	LogRegister(name string)
	LogObserverStart()
	LogObserverStop()
}

func NopLogger() Logger { return new(nopLogger) }

type nopLogger struct{}

func (nopLogger) LogRegister(_ string) {}

func (nopLogger) LogObserverStart() {}

func (nopLogger) LogObserverStop() {}

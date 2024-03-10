// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build dev

package main

import (
	"fmt"
	"github.com/go-pogo/env"
	"github.com/go-pogo/env/dotenv"
	"github.com/rs/zerolog"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func init() {
	zerolog.ErrorMarshalFunc = func(err error) interface{} {
		return fmt.Sprintf("%+v", err)
	}

	unmarshalEnv = func(v any) error {
		_, f, _, _ := runtime.Caller(1)

		ae, _ := dotenv.GetActiveEnvironmentOr(os.Args[1:], dotenv.Development)
		environ, err := dotenv.Read(filepath.Dir(f), ae).Environ()
		if err != nil {
			return err
		}

		if err = env.Load(environ); err != nil {
			return err
		}
		return env.NewDecoder(environ).Decode(v)
	}

	loggerOut = func() io.Writer {
		out := zerolog.NewConsoleWriter()
		out.TimeFormat = time.StampMilli
		return out
	}
}

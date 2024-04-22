// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build dev

package main

import (
	"fmt"
	"github.com/go-pogo/env"
	"github.com/go-pogo/env/dotenv"
	observerapp "github.com/roeldev/youless-observer/app/observer"
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

	unmarshalEnv = func(conf *observerapp.Config) error {
		_, dir, _, _ := runtime.Caller(1)
		dir = filepath.Dir(dir)

		ae, _ := dotenv.GetActiveEnvironmentOr(os.Args[1:], dotenv.Development)
		environ, err := dotenv.Read(dir, ae).Environ()
		if err != nil {
			return err
		}

		if err = env.Load(environ); err != nil {
			return err
		}
		if err = env.NewDecoder(environ).Decode(conf); err != nil {
			return err
		}

		prefixDir(dir, &conf.Server.TLS.CACertFile)
		prefixDir(dir, &conf.Server.TLS.CertFile)
		prefixDir(dir, &conf.Server.TLS.KeyFile)
		return nil
	}

	loggerOut = func() io.Writer {
		out := zerolog.NewConsoleWriter()
		out.TimeFormat = time.StampMilli
		return out
	}
}

func prefixDir(dir string, ptr *string) {
	if ptr == nil || *ptr == "" || filepath.IsAbs(*ptr) {
		return
	}
	*ptr = filepath.Join(dir, *ptr)
}

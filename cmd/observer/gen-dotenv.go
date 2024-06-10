// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"github.com/go-pogo/env/envfile"
	"github.com/go-pogo/errors"
	"github.com/roeldev/youless-observer/cmd/observer/observer-app"
	"path/filepath"
	"runtime"
)

func main() {
	_, dir, _, _ := runtime.Caller(0)
	dir = filepath.Dir(dir)
	errors.FatalOnErr(envfile.Generate(dir, ".env", observerapp.Config{}))
}

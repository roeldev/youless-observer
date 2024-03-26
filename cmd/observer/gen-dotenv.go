// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"github.com/go-pogo/env/envfile"
	"github.com/go-pogo/errors"
	"github.com/roeldev/youless-observer/app/observer"
	"log"
	"path/filepath"
	"runtime"
)

func main() {
	_, dir, _, _ := runtime.Caller(0)
	dir = filepath.Dir(dir)

	filename := filepath.Join(dir, ".env")
	if err := write(filename); err != nil {
		log.Printf("cannot write to %s: %s", filename, err)
	}
}

func write(filename string) (err error) {
	enc, err := envfile.Create(filename)
	if err != nil {
		return errors.WithStack(err)
	}

	enc.TakeValues = true
	defer errors.AppendFunc(&err, enc.Close)

	if err = enc.Encode(observerapp.Config{}); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

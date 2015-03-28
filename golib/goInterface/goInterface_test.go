// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package goInterface

import (
	//	"fmt"
	"testing"
)

func TestInitAppFilesDir(t *testing.T) {
	if err := InitAppFilesDir("files"); err != nil {
		t.Error(err)
	}
}

func TestStopNotStarted(t *testing.T) {
	//##NB; we did not start the provider
	if Stop() == nil {
		t.Error("provider was not started, Stop should return an error")
	}
}

func TestStart(t *testing.T) {
	if err := Start(); err != nil {
		t.Error(err)
	}
}

func TestStop(t *testing.T) {
	if err := Stop(); err != nil {
		t.Error(err)
	}
}

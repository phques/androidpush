// AndroidPush project
// Copyright 2015 Philippe Quesnel
// Licensed under the Academic Free License version 3.0

package goInterface

import (
	//	"fmt"
	"testing"
)

func TestStopNotStarted(t *testing.T) {
	//##NB; we did not start the provider
	if Stop() == nil {
		t.Error("provider was not started, Stop should return an error")
	}
}

func TestStart(t *testing.T) {
	// init before start !
	param := NewInitParam()
	param.AppFilesDir = "."
	if err := Init(param); err != nil {
		t.Error("failed to init:", err)
	}

	if err := Start(); err != nil {
		t.Error("failed to Start:", err)
	}
}

func TestStop(t *testing.T) {
	if err := Stop(); err != nil {
		t.Error("failed to Stop:", err)
	}
}

//go:build !windows

package watcher

import (
	"errors"
)

type StubWatcher struct{}

func New(callbacks EventCallbacks) Watcher {
	return &StubWatcher{}
}

func (w *StubWatcher) Start() error {
	return errors.New("win32 event watcher is only supported on Windows")
}

func (w *StubWatcher) Stop() {}

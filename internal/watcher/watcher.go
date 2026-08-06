package watcher

// EventCallbacks defines the handler functions for system events.
type EventCallbacks struct {
	OnSleep    func()
	OnWake     func()
	OnLock     func()
	OnUnlock   func()
	OnShutdown func()
}

// Watcher interface for system event listeners.
type Watcher interface {
	Start() error
	Stop()
}

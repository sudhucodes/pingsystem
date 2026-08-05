package telegram

import (
	"testing"
)

func TestEscapeMarkdown(t *testing.T) {
	input := "Hello_World! [v1.0] (Test-Build)"
	expected := `Hello\_World\! \[v1\.0\] \(Test\-Build\)`
	actual := escapeMarkdown(input)
	if actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestGetEventIcon(t *testing.T) {
	if getEventIcon(EventStartup) != "🚀" {
		t.Errorf("expected 🚀 for Startup event")
	}
	if getEventIcon(EventSleep) != "🌙" {
		t.Errorf("expected 🌙 for Sleep event")
	}
	if getEventIcon(EventWake) != "☀️" {
		t.Errorf("expected ☀️ for Wake event")
	}
	if getEventIcon(EventShutdown) != "🛑" {
		t.Errorf("expected 🛑 for Shutdown event")
	}
}

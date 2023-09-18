package logging

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestInit_WithCorrectLoggingLevels(t *testing.T) {
	var buf bytes.Buffer

	log_levels := []string{"debug", "info", "warn", "error"}
	expected_logs := [][]string{
		{"debug_line", "info_line", "warn_line", "error_line"},
		{"info_line", "warn_line", "error_line"},
		{"warn_line", "error_line"},
		{"error_line"}}
	unexpected_logs := [][]string{{},
		{"debug_line"},
		{"debug_line", "info_line"},
		{"debug_line", "info_line", "warn_line"}}

	for i, log_level := range log_levels {
		buf.Reset()
		ed := Init(log_level, &buf)
		if ed != nil {
			t.Error(ed)
		}
		slog.Debug("debug_line")
		slog.Info("info_line")
		slog.Warn("warn_line")
		slog.Error("error_line")

		for _, unexpected_log := range unexpected_logs[i] {
			if bytes.Contains(buf.Bytes(), []byte(unexpected_log)) {
				t.Error("unexpected log found: " + unexpected_log)
			}
		}
		for _, expected_log := range expected_logs[i] {
			if !bytes.Contains(buf.Bytes(), []byte(expected_log)) {
				t.Error("expected log not found: " + expected_log)
			}
		}
	}
}

func TestInit_WithIncorrectLoggingLevel(t *testing.T) {
	var buf bytes.Buffer

	ed := Init("wrong_log_level", &buf)
	if ed == nil {
		t.Error("expected error")
	}
}

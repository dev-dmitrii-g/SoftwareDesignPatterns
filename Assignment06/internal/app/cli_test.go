package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"testing"
)

func TestCLIReports(t *testing.T) {
	for _, report := range []string{"quality", "alarms", "filter", "json", "latency"} {
		t.Run(report, func(t *testing.T) {
			var output bytes.Buffer
			args := []string{"-samples", "120", "-target", "patient-1", "-report", report}
			if err := RunCLI(context.Background(), args, &output, io.Discard); err != nil {
				t.Fatal(err)
			}
			if !json.Valid(output.Bytes()) {
				t.Fatalf("invalid JSON: %s", output.String())
			}
			if report == "json" && (bytes.Contains(output.Bytes(), []byte("patient-2")) || !bytes.Contains(output.Bytes(), []byte("patient-1"))) {
				t.Fatalf("wrong patient subtree: %s", output.String())
			}
		})
	}
}

func TestCLIRejectsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{
		{"-workers", "0"},
		{"-workers", "9"},
		{"-samples", "0"},
		{"-report", "missing"},
		{"-target", "missing"},
		{"-unknown"},
		{"unexpected"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			if err := RunCLI(context.Background(), args, io.Discard, io.Discard); err == nil {
				t.Fatalf("accepted invalid arguments: %v", args)
			}
		})
	}
}

func TestCLIPreviewStopsAfterFiveSamples(t *testing.T) {
	var progress bytes.Buffer
	if err := RunCLI(context.Background(), []string{"-samples", "120"}, io.Discard, &progress); err != nil {
		t.Fatal(err)
	}
	previewLines := 0
	hasTotal := false
	for _, line := range strings.Split(progress.String(), "\n") {
		if strings.Contains(line, " value=") && strings.Contains(line, " quality=") {
			previewLines++
		}
		if strings.HasPrefix(line, "samples=120 ") {
			hasTotal = true
		}
	}
	if previewLines != 5 || !hasTotal {
		t.Fatalf("preview lines = %d, complete summary = %t\n%s", previewLines, hasTotal, progress.String())
	}
}

package tests

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/protocol"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/lsp/source"
	"github.com/govim/govim/cmd/govim/internal/golang_org_x_tools/span"
)

// DiffDiagnostics prints the diff between expected and actual diagnostics test
// results.
func DiffDiagnostics(uri span.URI, want, got []source.Diagnostic) string {
	source.SortDiagnostics(want)
	source.SortDiagnostics(got)

	if len(got) != len(want) {
		return summarizeDiagnostics(-1, uri, want, got, "different lengths got %v want %v", len(got), len(want))
	}
	for i, w := range want {
		g := got[i]
		if w.Message != g.Message {
			return summarizeDiagnostics(i, uri, want, got, "incorrect Message got %v want %v", g.Message, w.Message)
		}
		if w.Severity != g.Severity {
			return summarizeDiagnostics(i, uri, want, got, "incorrect Severity got %v want %v", g.Severity, w.Severity)
		}
		if w.Source != g.Source {
			return summarizeDiagnostics(i, uri, want, got, "incorrect Source got %v want %v", g.Source, w.Source)
		}
		// Don't check the range on the badimport test.
		if strings.Contains(uri.Filename(), "badimport") {
			continue
		}
		if protocol.ComparePosition(w.Range.Start, g.Range.Start) != 0 {
			return summarizeDiagnostics(i, uri, want, got, "incorrect Start got %v want %v", g.Range.Start, w.Range.Start)
		}
		if !protocol.IsPoint(g.Range) { // Accept any 'want' range if the diagnostic returns a zero-length range.
			if protocol.ComparePosition(w.Range.End, g.Range.End) != 0 {
				return summarizeDiagnostics(i, uri, want, got, "incorrect End got %v want %v", g.Range.End, w.Range.End)
			}
		}
	}
	return ""
}

func summarizeDiagnostics(i int, uri span.URI, want []source.Diagnostic, got []source.Diagnostic, reason string, args ...interface{}) string {
	msg := &bytes.Buffer{}
	fmt.Fprint(msg, "diagnostics failed")
	if i >= 0 {
		fmt.Fprintf(msg, " at %d", i)
	}
	fmt.Fprint(msg, " because of ")
	fmt.Fprintf(msg, reason, args...)
	fmt.Fprint(msg, ":\nexpected:\n")
	for _, d := range want {
		fmt.Fprintf(msg, "  %s:%v: %s\n", uri, d.Range, d.Message)
	}
	fmt.Fprintf(msg, "got:\n")
	for _, d := range got {
		fmt.Fprintf(msg, "  %s:%v: %s\n", uri, d.Range, d.Message)
	}
	return msg.String()
}

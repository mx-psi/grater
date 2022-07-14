package scraper

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependents(t *testing.T) {
	tests := []struct {
		inputPath  string
		outputPath string
	}{
		{
			inputPath:  "testdata/webpages/go.opentelemetry.io_collector_client_0.55.0.html",
			outputPath: "testdata/output/go.opentelemetry.io_collector_client_0.55.0-dependents.json",
		},
		{
			inputPath:  "testdata/webpages/go.opentelemetry.io_collector_component _0.55.html",
			outputPath: "testdata/output/go.opentelemetry.io_collector_component_0.55.0-dependents.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.inputPath, func(t *testing.T) {
			in, err := os.Open(tt.inputPath)
			require.NoError(t, err)

			dependents, err := dependentsFromReader(in)
			require.NoError(t, err)

			var expected []ModuleDep
			out, err := os.ReadFile(tt.outputPath)
			require.NoError(t, err)
			err = json.Unmarshal(out, &expected)
			require.NoError(t, err)

			assert.ElementsMatch(t, expected, dependents)
		})
	}
}

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		path        string
		cfg         Config
		err         string
		validateErr string
	}{
		{
			path: "testdata/invalid.yaml",
			err:  "yaml: unmarshal errors:\n  line 1: field unknown_key not found in type config.Config",
		},
		{
			path: "testdata/module.yaml",
			cfg: Config{Modules: []ModuleConfig{
				{
					Path:          "go.opentelemetry.io/collector/component",
					BaseVersion:   "v0.55.0",
					TargetVersion: "v0.56.0",
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			cfg, err := Load(tt.path)
			if err != nil || tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Equal(t, tt.cfg, cfg)
				if err := cfg.Validate(); err != nil || tt.validateErr != "" {
					assert.EqualError(t, err, tt.validateErr)
				}
			}
		})
	}
}

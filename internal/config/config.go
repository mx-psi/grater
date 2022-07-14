package config

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// validatePath of configuration.
func validatePath(path string) error {
	validExtensions := []string{".yaml", ".yml"}
	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			return nil
		}
	}
	return fmt.Errorf("%q does not have a valid extension (%v)", path, validExtensions)
}

// LoadAndValidate Config from file.
func LoadAndValidate(path string) (cfg Config, err error) {
	cfg, err = Load(path)
	if err != nil {
		return
	}

	if err = cfg.Validate(); err != nil {
		return
	}

	return
}

// Load Config from file.
func Load(path string) (cfg Config, err error) {
	if err = validatePath(path); err != nil {
		return
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.UnmarshalStrict(bytes, &cfg)
	return
}

// Config of grater.
type Config struct {
	// Modules to be tested.
	Modules []ModuleConfig `yaml:"modules"`
}

// Validate the configuration.
func (cfg *Config) Validate() error {
	if nmod := len(cfg.Modules); nmod != 1 {
		return fmt.Errorf("found %d modules, exactly 1 is currently supported", nmod)
	}
	return nil
}

// PathConfig of a Go module.
type PathConfig string

// VersionConfig of a Go module.
type VersionConfig string

// ModuleConfig is a config representation of a Go module.
type ModuleConfig struct {
	// Path of the Go module.
	Path PathConfig `yaml:"path"`
	// BaseVersion of the Go module.
	BaseVersion VersionConfig `yaml:"base_version"`
	// TargetVersion of the Go module.
	TargetVersion VersionConfig `yaml:"target_version"`
}

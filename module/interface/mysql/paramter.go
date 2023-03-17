package mysql

import (
	"github.com/hashicorp/go-version"
)

type Mode interface {
	String() string
}

type Parameter interface {
	// GetConfig gets the configuration of Parameter
	GetConfig(v *version.Version, mode ...Mode) ([]byte, error)
	// WriteConfig writes the configuration to the specified file
	WriteConfig(configPath string, data []byte) error
	// Marshal marshals the Parameter to json bytes
	Marshal() ([]byte, error)
	// Unmarshal unmarshals the specified json bytes to Parameter
	Unmarshal(b []byte) error
}

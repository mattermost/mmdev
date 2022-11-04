package model

import "fmt"

type VersionConfig struct {
	MinVersion string
	MaxVersion string
}

type Config struct {
	NodeJS VersionConfig
	NVM    VersionConfig
}

func DefaultConfig() *Config {
	return &Config{
		NodeJS: VersionConfig{
			MinVersion: "18",
		},
		NVM: VersionConfig{
			MinVersion: "0.39.1",
		},
	}
}

func (c VersionConfig) String() string {
	if c.MinVersion != "" && c.MaxVersion != "" {
		return fmt.Sprintf("between %s and %s", c.MinVersion, c.MaxVersion)
	} else if c.MinVersion != "" {
		return c.MinVersion
	} else {
		return c.MaxVersion
	}
}

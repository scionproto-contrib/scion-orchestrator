package conf

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type InstallationConfig struct {
	ConfigPath   string `toml:"config_path"`
	DatabasePath string `toml:"database_path"`
	LogPath      string `toml:"log_path"`

	BinaryPath  string `toml:"binary_path"`
	ServicePath string `toml:"service_path"`

	InstalledAt string `toml:"installed_at"`
	LastUpdate  string `toml:"last_update"`

	ScionVersion string `toml:"scion_version"`
}

func NewInstallationConfig() *InstallationConfig {
	return &InstallationConfig{}
}

func LoadInstallationConfig(path string) (*InstallationConfig, error) {
	c := NewInstallationConfig()

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(content, c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *InstallationConfig) Save(path string) error {
	content, err := toml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, 0644)
}

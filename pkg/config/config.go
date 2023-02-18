package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/phayes/permbits"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ConfigPath is the path to our configuration file on disk
var ConfigPath string

type Config struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	SessionState string `yaml:"session_state"`
	RedirectURL  string `yaml:"redirect_url"`
}

// NewConfig creates a new, empty, config
func NewConfig() *Config {
	return &Config{}
}

// Save saves a configuration
func (c *Config) Save() error {
	yamlByte, err := yaml.Marshal(c)
	if err != nil {
		return errors.Wrap(err, "Unable to serialize configuration")
	}

	err = ioutil.WriteFile(ConfigPath, yamlByte, 0600)
	if err != nil {
		return errors.Wrap(err, "Failed to write configuration")
	}

	return nil
}

// Read reads the config file and creates a empty config file if could not find a config file at the given path
func Read() (*Config, error) {
	err := EnsureConfig()
	if err != nil {
		return NewConfig(), errors.Wrap(err, "Failed to read config")
	}

	data, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return NewConfig(), errors.Wrap(err, "Failed to read config")
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return NewConfig(), errors.Wrap(err, "Failed to parse config")
	}

	return config, nil
}

// EnsureConfig ensures that the config is there and exists using the correct permissions.
func EnsureConfig() error {
	// Create config file if not exists
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		path, _ := path.Split(ConfigPath)
		os.MkdirAll(path, 0600)

		config := NewConfig()
		err = config.Save()
		if err != nil {
			return errors.Wrap(err, "failed to write empty config")
		}
	}

	return checkPermissions(ConfigPath)
}

func checkPermissions(configPath string) error {
	permissions, err := permbits.Stat(configPath)
	if err != nil {
		return errors.Wrap(err, "failed to check file permissions for config")
	}

	if permissions.GroupRead() || permissions.GroupWrite() || permissions.GroupExecute() ||
		permissions.OtherRead() || permissions.OtherWrite() || permissions.OtherExecute() {
		fmt.Println("")
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("!! Security Alert")
		fmt.Println("!! config is world readable!")
		fmt.Println("!! Since contains OIDC client secrets, this could cause serious security issues!")
		fmt.Println("!! to get rid of this message, run")
		fmt.Printf("!! $ chmod -R 600 %s\n", configPath)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("")
	} else if permissions.UserExecute() {
		fmt.Println("")
		fmt.Println("!! Configuration anomaly detected")
		fmt.Println("!! config should not be executable")
		fmt.Println("!! to get rid of this message, run")
		fmt.Printf("!! $ chmod -R 600 %s\n", configPath)
		fmt.Println("")
	}

	return nil
}

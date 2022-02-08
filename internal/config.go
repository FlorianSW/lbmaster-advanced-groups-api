package internal

import (
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"os"
)

type Config struct {
	Port                     int    `json:"port"`
	AdvancedGroupsConfigPath string `json:"advanced_groups_config_path"`
	ApiKey                   string `json:"api_key"`
}

func NewConfig(path string, logger lager.Logger) (Config, error) {
	config, err := readConfig(path, logger)
	if err != nil {
		return config, err
	}

	return config, persistConfig(path, config)
}

func readConfig(path string, logger lager.Logger) (Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Info("create-config")
		config = Config{
			Port: 8080,
		}
	} else {
		logger.Info("read-existing-config")
		c, err := os.ReadFile(path)
		if err != nil {
			return Config{}, err
		}
		err = json.Unmarshal(c, &config)
		if err != nil {
			return Config{}, err
		}
	}
	return config, nil
}

func persistConfig(path string, config Config) error {
	c, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, c, 0655)
}

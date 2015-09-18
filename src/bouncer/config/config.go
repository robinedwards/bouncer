package config

import (
	"bouncer/experiment"
	"encoding/json"
	"bouncer/feature"
	"fmt"
)

type Config struct {
	Experiments []experiment.Experiment
	Groups      []Group
	Features	[]feature.Feature
}

type Group struct {
	Name string
	Uids []string
}

// todo: doc
func LoadConfig(jsonConfig string) (Config, error) {
	if len(jsonConfig) == 0 {
		return Config{}, nil
	}

	var config Config
	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		fmt.Println(err)
	}

	return config, err
}

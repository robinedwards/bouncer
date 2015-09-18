package config

import (
	"bouncer/experiment"
	"encoding/json"
	"fmt"
)

type Config struct {
	Experiments []experiment.Experiment
	Groups      []Group
}

type Group struct {
	Name string
	Uids []string
}

// todo: doc
func LoadConfig(jsonConfig string) Config {
	if len(jsonConfig) == 0 {
		return Config{}
	}

	var config Config
	err := json.Unmarshal([]byte(jsonConfig), &config)
	if err != nil {
		// todo log / raise error
		fmt.Println("error:", err)
	}

	return config
}

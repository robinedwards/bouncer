package config

import (
	"bouncer/experiment"
	"encoding/json"
	"bouncer/feature"
	"fmt"
	"io/ioutil"
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

func LoadConfig(jsonConfig string) (Config, error) {
	if len(jsonConfig) == 0 {
		return Config{}, nil
	}

	var config Config
	err := json.Unmarshal([]byte(jsonConfig), &config)

	return config, err
}

func LoadConfigFile(filename string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error loading file: %v\n", err)
		return config, err
    }

    merr := json.Unmarshal(file, &config)
    return config, merr
}

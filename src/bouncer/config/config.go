package config

import (
	"bouncer/experiment"
	"bouncer/feature"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Experiments []experiment.Experiment
	Groups      []Group
	Features    []feature.Feature
}

type Group struct {
	Name string
	Uids []string
}

// Take group configurations and wire them into features and experiments
func patchInGroups(config Config) error {

	groupMap := make(map[string][]string)

	for _, group := range config.Groups {
		groupMap[group.Name] = group.Uids
	}

	for _, feature := range config.Features {
		err := feature.SetupGroups(groupMap)
		if err != nil {
			return err
		}
	}

	for _, experiment := range config.Experiments {
		err := experiment.SetupGroups(groupMap)
		if err != nil {
			return err
		}
	}

	return nil
}

func LoadConfig(jsonConfig string) (Config, error) {
	if len(jsonConfig) == 0 {
		return Config{}, nil
	}

	var config Config
	merr := json.Unmarshal([]byte(jsonConfig), &config)
	if merr != nil {
		return config, merr
	}

	err := patchInGroups(config)

	return config, err
}

func LoadConfigFile(filename string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("%v\n", err)
		return config, err
	}

	merr := json.Unmarshal(file, &config)
	if merr != nil {
		return config, merr
	}

	perr := patchInGroups(config)

	return config, perr
}

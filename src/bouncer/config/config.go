package config

import (
	"bouncer/experiment"
	"bouncer/feature"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Experiments   []*experiment.Experiment `json:"experiments"`
	Groups        []Group                  `json:"groups"`
	Features      []*feature.Feature       `json:"features"`
	FeatureMap    map[string]*feature.Feature
	ExperimentMap map[string]*experiment.Experiment
}

type Group struct {
	Name string   `json:"name"`
	Uids []string `json:"uids"`
}

// Take group configurations and wire them into features and experiments
func setupGroupsAndRings(config *Config) error {

	groupMap := make(map[string][]string)

	for _, group := range config.Groups {
		groupMap[group.Name] = group.Uids
	}

	for _, feature := range config.Features {
		err := feature.SetupGroups(groupMap)
		feature.SetupRing()
		if err != nil {
			return err
		}
	}

	for _, experiment := range config.Experiments {
		err := experiment.SetupGroups(groupMap)
		if err != nil {
			return err
		}
		experiment.SetupRing()
	}

	return nil
}

// Create feature and experiment mappings
func setupMappings(config *Config) {
	config.FeatureMap = make(map[string]*feature.Feature)

	for _, feature := range config.Features {
		config.FeatureMap[feature.Name] = feature
	}

	config.ExperimentMap = make(map[string]*experiment.Experiment)

	for _, experiment := range config.Experiments {
		config.ExperimentMap[experiment.Name] = experiment
	}
}

func InitConfig(config *Config) error {
	err := setupGroupsAndRings(config)
	setupMappings(config)
	return err
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

	err := InitConfig(&config)

	return config, err
}

func LoadConfigFile(filename string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Errorf("%v\n", err)
		return config, err
	}

	merr := json.Unmarshal(file, &config)
	if merr != nil {
		return config, merr
	}

	perr := InitConfig(&config)

	return config, perr
}

package config_test
import (
	"testing"
	"bouncer/config"
	"bouncer/experiment"
	"reflect"
)

var exampleJsonConfig = `
      {
        "groups": [
          {
            "name": "admins",
            "uids": [
              "8",
              "32"
            ]
          },
          {
            "name": "testers",
            "uids": [
              "8",
              "32"
            ]
          }
        ],
        "features": {
          "audio_mode": {
            "groups": {
              "admins": 1,
              "users": 0
            },
            "enabled": 0.5
          }
        },
        "experiments": [
          {
            "name": "progress_bar",
            "groups": [
              {
                "admins": "green",
                "testers": "red"
              }
            ],
            "alternatives": [
              {
                "name": "green",
                "weight": 1
              },
              {
                "name": "red",
                "weight": 9
              }
            ]
          }
        ]
      }`

func validateEmptyConfig(configObj config.Config) bool {
	return len(configObj.Experiments) == 0 && len(configObj.Groups) == 0
}

func TestEmptyJson(t *testing.T) {
	if !validateEmptyConfig(config.LoadConfig(``)) {
		t.Error("empty config-string should return a default config")
	}
}

func TestInvalidJson(t *testing.T) {
	if !validateEmptyConfig(config.LoadConfig(`{`)) {
		t.Error("empty config-string should return a default config")
	}
}

func TestValidJson(t *testing.T) {
	testConfig := config.LoadConfig(exampleJsonConfig)
	if len(testConfig.Experiments) != 1 {
		t.Error("Example config should contain 1 experiments. Found:", len(testConfig.Experiments))
	}

	if len(testConfig.Groups) != 2 {
		t.Error("Example config should contain 2 groups. Found:", len(testConfig.Groups))
	}

	// todo danielz: fix this test
	if reflect.DeepEqual(testConfig.Experiments[0],
		experiment.NewExperiment("progress_bar", experiment.Alternative{"green", 1}, experiment.Alternative{"red", 1})) {
		t.Error("Wrong experiment for example config")
	}

}
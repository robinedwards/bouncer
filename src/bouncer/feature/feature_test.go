package feature_test

import (
	"bouncer/feature"
	"strconv"
	"testing"
)

func TestFeature(t *testing.T) {

	f := feature.NewFeature("video", 1, make(map[string]int))

	if !f.IsEnabled("23424235") {
		t.Errorf("Should be disabled")
	}
}

func TestFeatureDisabled(t *testing.T) {

	f := feature.NewFeature("video", 0, make(map[string]int))

	if f.IsEnabled("3452355") {
		t.Errorf("Should be disabled")
	}
}

func TestFeaturePartEnabled(t *testing.T) {
	f := feature.NewFeature("video", 0.5, make(map[string]int))
	enabled := 0
	disabled := 0
	for i := 0; i < 10; i++ {
		if f.IsEnabled(strconv.Itoa(i)) {
			enabled += 1
		} else {
			disabled += 1
		}
	}

	if enabled == 0 || disabled == 0 {
		t.Errorf("Should be mix of disabled and enabled")
	}
}

func TestFeatureWithGroupMapping(t *testing.T)  {
	groupConfig := map[string]int{
		"group_a": 1,
		"group_b": 0,
	}

	groupMapping := map[string][]string{
		"group_a": {"1", "2"},
		"group_b": {"3", "4"},
	}

	f := feature.NewFeature("video", 0, groupConfig)
	f.SetupGroups(groupMapping)

	if !f.IsEnabled("1") {
		t.Errorf("uid 1 is in group a should be enabled")
	}

	if f.IsEnabled("3") {
		t.Errorf("uid 3 is in group b should be disabled")
	}
}
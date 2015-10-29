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

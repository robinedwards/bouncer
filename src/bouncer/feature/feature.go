package feature

import (
	"github.com/serialx/hashring"
)

type Feature struct {
	Name    string
	enabled float32
	ring    *hashring.HashRing
}

func NewFeature(name string, enabled float32) Feature {
	if enabled > 1 || enabled < 0 {
		panic("enabled must be between > 0 and <= 1")
	}

	var ring *hashring.HashRing

	if enabled > 0 && enabled < 1 {
		weights := make(map[string]int)

		weights["enabled"] = int(enabled * 100)
		weights["disabled"] = int((1 - enabled) * 100)
		ring = hashring.NewWithWeights(weights)
	}

	return Feature{name, enabled, ring}
}

func (f Feature) IsEnabled(uid string) bool {
	if f.ring == nil {
		return f.enabled == 1
	}

	r, _ := f.ring.GetNode(uid)

	return r == "enabled"
}

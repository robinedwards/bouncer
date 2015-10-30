package feature

import (
	"github.com/serialx/hashring"
)

type Feature struct {
	Name       string
	Enabled    float32
	ring       *hashring.HashRing
	groups     map[string]int
	uidMapping map[string]int
}

func NewFeature(name string, enabled float32, groups map[string]int) Feature {
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

	return Feature{name, enabled, ring, groups, make(map[string]int)}
}

func (f Feature) SetupGroups(groups map[string][]string) error {
	for groupName, status := range f.groups {
		groupUids := groups[groupName]

		for _, uid := range groupUids {
			f.uidMapping[uid] = status
		}
	}

	return nil
}

func (f Feature) IsEnabled(uid string) bool {
	// has uid been overriden in a group
	if status, ok := f.uidMapping[uid]; ok {
		return status == 1
	}

	if f.ring == nil {
		return f.Enabled == 1
	}

	r, _ := f.ring.GetNode(uid)

	return r == "enabled"
}

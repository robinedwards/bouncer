package feature

import (
	"fmt"
	"github.com/serialx/hashring"
)

type Feature struct {
	Name       string  `json:"name"`
	Enabled    float32 `json:"enabled"`
	ring       *hashring.HashRing
	groups     map[string]int
	uidMapping map[string]int
}

func NewFeature(name string, enabled float32, groups map[string]int) Feature {
	if enabled > 1 || enabled < 0 {
		panic("enabled must be between > 0 and <= 1")
	}

	f := Feature{name, enabled, nil, groups, make(map[string]int)}
	f.SetupRing()
	return f
}

func (f *Feature) SetupRing() {
	if f.Enabled > 0.0 && f.Enabled < 1.0 {
		weights := make(map[string]int)
		weights["enabled"] = int(f.Enabled * 100)
		weights["disabled"] = int((1 - f.Enabled) * 100)
		f.ring = hashring.NewWithWeights(weights)
	}
}

func (f *Feature) SetupGroups(groups map[string][]string) error {
	for groupName, status := range f.groups {
		groupUids := groups[groupName]

		for _, uid := range groupUids {
			f.uidMapping[uid] = status
		}
	}

	return nil
}

func (f *Feature) IsEnabled(uid string) bool {
	// has uid been overriden in a group
	if status, ok := f.uidMapping[uid]; ok {
		return status == 1
	}

	if f.ring == nil {
		return f.Enabled == 1
	}

	r, err := f.ring.GetNode(uid)
	if !err {
		fmt.Errorf("ERROR: getting node for feature, disabling by default")
		return false
	}
	return r == "enabled"
}

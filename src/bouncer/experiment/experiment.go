package experiment

import (
	"fmt"
	"github.com/serialx/hashring"
)

type Experiment struct {
	Name         string        `json:"name"`
	Alternatives []Alternative `json:"alternatives"`
	ring         *hashring.HashRing
	groups       map[string]string
	uidMapping   map[string]string
}

type Alternative struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
}

func NewExperiment(name string, groups map[string]string, alternatives ...Alternative) Experiment {
	exp := Experiment{name, alternatives, nil, groups, make(map[string]string)}
	exp.SetupRing()
	return exp
}

func (exp *Experiment) GetAlternative(uid string) string {
	// if we have an entry from a group
	if groupAlt, ok := exp.uidMapping[uid]; ok {
		return groupAlt
	}
	alt, err := exp.ring.GetNode(uid)

	if !err {
		fmt.Errorf("ERROR: getting alternate, defaulting to first")
		return exp.Alternatives[0].Name
	}

	return alt
}

func (exp *Experiment) SetupGroups(groups map[string][]string) error {
	for groupName, alternativeName := range exp.groups {
		groupUids := groups[groupName]

		found := false

		for _, alt := range exp.Alternatives {
			if alt.Name == alternativeName {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("ERROR: Alternative %s from group %s not in experiment %s",
				alternativeName, groupName, exp.Name)
		}

		for _, uid := range groupUids {
			exp.uidMapping[uid] = alternativeName
		}
	}

	return nil
}

func (exp *Experiment) SetupRing() {
	weights := make(map[string]int)

	for _, alt := range exp.Alternatives {
		weights[alt.Name] = alt.Weight
	}

	exp.ring = hashring.NewWithWeights(weights)
}

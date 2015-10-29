package experiment

import (
	"fmt"
	"github.com/serialx/hashring"
)

type Experiment struct {
	Name         string
	Alternatives []Alternative
	ring         hashring.HashRing
	groups       map[string]string
	uidMapping   map[string]string
}

type Alternative struct {
	Name   string
	Weight int
}

func NewExperiment(name string, groups map[string]string, alternatives ...Alternative) Experiment {
	weights := make(map[string]int)

	for _, alt := range alternatives {
		weights[alt.Name] = alt.Weight
	}

	ring := hashring.NewWithWeights(weights)

	return Experiment{name, alternatives, *ring, groups, make(map[string]string)}
}

func (exp *Experiment) GetAlternative(uid string) string {
	// if we have an entry from a group
	if groupAlt, ok := exp.uidMapping[uid]; ok {
		return groupAlt
	}

	alt, err := exp.ring.GetNode(uid)

	if !err {
		panic(err)
	}

	return alt
}

func (exp Experiment) SetupGroups(groups map[string][]string) error {
	for groupName, alternativeName := range exp.groups {
		groupUids := groups[groupName]

		found := false

		for _, alt := range exp.Alternatives {
			if alt.Name == alternativeName {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Alternative %s from group %s not in experiment %s",
				alternativeName, groupName, exp.Name)
		}

		for _, uid := range groupUids {
			exp.uidMapping[uid] = alternativeName
		}
	}

	return nil
}

func Participate(experiments []Experiment, uid string) map[string]string {
	r := make(map[string]string)

	for _, exp := range experiments {
		alt := exp.GetAlternative(uid)
		r[exp.Name] = alt
	}

	return r
}

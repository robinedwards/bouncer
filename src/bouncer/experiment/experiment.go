package experiment

import (
	"github.com/serialx/hashring"
	"strconv"
)

type Experiment struct {
	Name         string
	Alternatives []Alternative
	ring         hashring.HashRing
}

type Alternative struct {
	Name   string
	Weight int
}

func NewExperiment(name string, alternatives ...Alternative) Experiment {
	weights := make(map[string]int)

	for _, alt := range alternatives {
		weights[alt.Name] = alt.Weight
	}

	ring := hashring.NewWithWeights(weights)

	return Experiment{name, alternatives, *ring}
}

func (exp *Experiment) GetAlternative(uid int64) string {
	alt, err := exp.ring.GetNode(strconv.FormatInt(uid, 10))

	if !err {
		panic(err)
	}

	return alt
}

func Participate(experiments []Experiment, uid int64) map[string]string {
	r := make(map[string]string)

	for _, exp := range experiments {
		alt := exp.GetAlternative(uid)
		r[exp.Name] = alt
	}

	return r
}

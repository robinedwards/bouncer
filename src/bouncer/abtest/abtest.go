package abtest

import (
	"github.com/serialx/hashring"
	"strconv"
)

type ABTest struct {
	Name         string
	Alternatives []Alternative
	ring         hashring.HashRing
}

type Alternative struct {
	Name   string
	Weight int
}

func NewABTest(name string, alternatives ...Alternative) ABTest {
	weights := make(map[string]int)

	for _, alt := range alternatives {
		weights[alt.Name] = alt.Weight
	}

	ring := hashring.NewWithWeights(weights)

	return ABTest{name, alternatives, *ring}
}

func (test *ABTest) GetAlternative(uid int64) string {
	alt, err := test.ring.GetNode(strconv.FormatInt(uid, 10))

	if !err {
		panic(err)
	}

	return alt
}

func Participate(tests []ABTest, uid int64) map[string]string {
	r := make(map[string]string)

	for _, test := range tests {
		alt := test.GetAlternative(uid)
		r[test.Name] = alt
	}

	return r
}

package abtest_test

import "testing"
import "abtest"

func TestBasicAB(t *testing.T) {
	newTest := abtest.NewABTest("test",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1})

	if newTest.Name != "test" {
		t.Errorf("Incorrect test name")
	}
}

func TestGetAlternate(t *testing.T) {
	newTest := abtest.NewABTest("test",
		abtest.Alternative{Name: "a", Weight: 1}, abtest.Alternative{Name: "b", Weight: 1})

	alternate := newTest.GetAlternative(1)

	if !(alternate == "a" || alternate == "b") {
		t.Errorf("Invalid alternate returned %s", alternate)
	}

	alternate1 := newTest.GetAlternative(1)

	if alternate != alternate1 {
		t.Errorf("Inconsistent hash")
	}
}

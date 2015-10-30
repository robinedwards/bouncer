package experiment_test

import "testing"
import "bouncer/experiment"

func TestBasicAB(t *testing.T) {
	newTest := experiment.NewExperiment("test",
		make(map[string]string),
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1})

	if newTest.Name != "test" {
		t.Errorf("Incorrect test name")
	}
}

func TestGetAlternate(t *testing.T) {
	newTest := experiment.NewExperiment("test", make(map[string]string),
		experiment.Alternative{Name: "a", Weight: 1}, experiment.Alternative{Name: "b", Weight: 1})

	alternate := newTest.GetAlternative("1")

	if !(alternate == "a" || alternate == "b") {
		t.Errorf("Invalid alternate returned %s", alternate)
	}

	alternate1 := newTest.GetAlternative("1")

	if alternate != alternate1 {
		t.Errorf("Inconsistent hash")
	}
}

func TestGroupMemberGetsCorrectAlternate(t *testing.T) {
	groupConfig := map[string]string{
		"group_a": "yellow_button",
		"group_b": "red_button",
	}

	groupMapping := map[string][]string{
		"group_a": {"1", "2"},
		"group_b": {"3", "4"},
	}

	newTest := experiment.NewExperiment("test",
		groupConfig,
		experiment.Alternative{Name: "yellow_button", Weight: 1},
		experiment.Alternative{Name: "red_button", Weight: 1})

	newTest.SetupGroups(groupMapping)

	if newTest.GetAlternative("1") != "yellow_button" {
		t.Errorf("Expecting yellow button for uid 1 in group a")
	}

	if newTest.GetAlternative("3") != "red_button" {
		t.Errorf("Expecting red button for uid 3 in group b")
	}
}

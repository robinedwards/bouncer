package handlers_test

import (
	"bouncer/config"
	"bouncer/experiment"
	"bouncer/feature"
	"bouncer/handlers"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockConfig() config.Config {
	cfg := config.Config{}
	cfg.Experiments = append(cfg.Experiments, experiment.NewExperiment("test1",
		make(map[string]string),
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

	cfg.Experiments = append(cfg.Experiments, experiment.NewExperiment("test2",
		make(map[string]string),
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

	cfg.Features = append(cfg.Features, feature.NewFeature("scrolling", 1, make(map[string]int)))

	return cfg
}

func checkValidResponse(code int, w *httptest.ResponseRecorder, t *testing.T) {
	if w.Code != code {
		t.Errorf("Expected %s got %s", code, w.Code)
	}
	_, err := json.Marshal(w.Body)

	if err != nil {
		t.Errorf("Response didn't return valid json")
	}
}

func TestListExperiments(t *testing.T) {
	mockCfg := mockConfig()

	h := handlers.ListExperiments(mockCfg)
	req, _ := http.NewRequest("GET", "/experiments/", nil)
	w := httptest.NewRecorder()

	h(w, req)

	checkValidResponse(http.StatusOK, w, t)
}

func TestListFeatures(t *testing.T) {
	mockCfg := mockConfig()

	h := handlers.ListFeatures(mockCfg)
	req, _ := http.NewRequest("GET", "/features/", nil)
	w := httptest.NewRecorder()

	h(w, req)
	checkValidResponse(http.StatusOK, w, t)
}

func TestListGroups(t *testing.T) {
	mockCfg := mockConfig()

	h := handlers.ListGroups(mockCfg)
	req, _ := http.NewRequest("GET", "/groups/", nil)
	w := httptest.NewRecorder()

	h(w, req)
	checkValidResponse(http.StatusOK, w, t)
}

func TestParticipate(t *testing.T) {
	mockCfg := mockConfig()

	h := handlers.Participate(mockCfg)
	req, _ := http.NewRequest("GET", "/participate/?uid=1", nil)
	w := httptest.NewRecorder()

	h(w, req)
	checkValidResponse(http.StatusOK, w, t)

	resp := new(handlers.ParticipateResponse)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading body %v", err)
	}

	merr := json.Unmarshal(body, resp)
	if merr != nil {
		t.Errorf("Error unmarshalling %v", merr)
	}

	if _, ok := resp.Experiments["test1"]; !ok {
		t.Errorf("Couldn't find test1 in the response", merr)
	}

	if _, ok := resp.Features["scrolling"]; !ok {
		t.Errorf("Couldn't find scrolling in the response", merr)
	}
}

func TestBadParticipate(t *testing.T) {
	mockCfg := mockConfig()

	h := handlers.Participate(mockCfg)
	req, _ := http.NewRequest("GET", "/participate/?n=f", nil)
	w := httptest.NewRecorder()

	h(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Participate page didn't return %v", http.StatusBadRequest)
	}
}

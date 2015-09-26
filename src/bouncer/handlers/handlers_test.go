package handlers_test

import (
	"bouncer/experiment"
	"bouncer/config"
	"bouncer/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
)

func mockConfig() config.Config{
	cfg := config.Config{}
	cfg.Experiments = append(cfg.Experiments, experiment.NewExperiment("test1",
				experiment.Alternative{Name: "a", Weight: 1},
				experiment.Alternative{Name: "b", Weight: 1}))

	cfg.Experiments = append(cfg.Experiments, experiment.NewExperiment("test2",
			experiment.Alternative{Name: "a", Weight: 1},
			experiment.Alternative{Name: "b", Weight: 1}))

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

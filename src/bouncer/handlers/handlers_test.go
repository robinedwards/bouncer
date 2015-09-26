package handlers_test

import (
	"bouncer/experiment"
	"bouncer/config"
	"bouncer/handlers"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestListExperiments(t *testing.T) {
	mockCfg := mockConfig()

	homeHandle := handlers.ListExperiments(mockCfg)
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	homeHandle(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestParticipate(t *testing.T) {
	mockCfg := mockConfig()

	homeHandle := handlers.Participate(mockCfg)
	req, _ := http.NewRequest("GET", "/participate/?uid=1", nil)
	w := httptest.NewRecorder()

	homeHandle(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Participate page didn't return %v", http.StatusOK)
	}
}

func TestBadParticipate(t *testing.T) {
	mockCfg := mockConfig()

	homeHandle := handlers.Participate(mockCfg)
	req, _ := http.NewRequest("GET", "/participate/?n=f", nil)
	w := httptest.NewRecorder()

	homeHandle(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Participate page didn't return %v", http.StatusBadRequest)
	}
}

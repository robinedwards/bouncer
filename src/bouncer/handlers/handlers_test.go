package handlers_test

import (
	"bouncer/config"
	"bouncer/experiment"
	"bouncer/feature"
	"bouncer/handlers"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockConfig() *config.Config {
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

	config.InitConfig(&cfg)

	return &cfg
}

func fakeLogger(interface{}) {}

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
	h := handlers.ListExperiments(mockConfig)
	req, _ := http.NewRequest("GET", "/experiments/", nil)
	w := httptest.NewRecorder()

	h(w, req)

	checkValidResponse(http.StatusOK, w, t)
}

func TestListFeatures(t *testing.T) {
	h := handlers.ListFeatures(mockConfig)
	req, _ := http.NewRequest("GET", "/features/", nil)
	w := httptest.NewRecorder()

	h(w, req)
	checkValidResponse(http.StatusOK, w, t)
}

func TestListGroups(t *testing.T) {
	h := handlers.ListGroups(mockConfig)
	req, _ := http.NewRequest("GET", "/groups/", nil)
	w := httptest.NewRecorder()

	h(w, req)
	checkValidResponse(http.StatusOK, w, t)
}

func makeParticipateRequest(req handlers.ParticipateRequest) httptest.ResponseRecorder {
	h := handlers.Participate(mockConfig, fakeLogger)
	body, _ := json.Marshal(req)

	r, _ := http.NewRequest("POST", "/participate/", bytes.NewReader(body))
	r.Header.Add("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, r)
	return *w
}

func checkParticipateResponse(w httptest.ResponseRecorder, t *testing.T) handlers.ParticipateResponse {
	checkValidResponse(http.StatusOK, &w, t)

	resp := new(handlers.ParticipateResponse)

	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Errorf("Error reading body %v", err)
	}

	merr := json.Unmarshal(body, resp)
	if merr != nil {
		t.Errorf("Error unmarshalling %v", merr)
	}
	return *resp
}

func TestParticipateSpecificTests(t *testing.T) {
	w := makeParticipateRequest(handlers.ParticipateRequest{
		Context:     handlers.Context{"1"},
		Experiments: map[string][]string{"test1": {"a", "b"}},
		Features:    map[string]float32{"scrolling": 1},
	})

	resp := checkParticipateResponse(w, t)

	if _, ok := resp.Experiments["test1"]; !ok {
		t.Errorf("Couldn't find test1 in the response")
	}

	if _, ok := resp.Features["scrolling"]; !ok {
		t.Errorf("Couldn't find scrolling in the response")
	}
}

func TestBasicParticipate(t *testing.T) {
	w := makeParticipateRequest(handlers.ParticipateRequest{
		Context:     handlers.Context{"1"},
		Experiments: map[string][]string{},
		Features:    map[string]float32{},
	})

	resp := checkParticipateResponse(w, t)

	if _, ok := resp.Experiments["test1"]; !ok {
		t.Errorf("Couldn't find test1 in the response")
	}

	if _, ok := resp.Features["scrolling"]; !ok {
		t.Errorf("Couldn't find scrolling in the response")
	}
}

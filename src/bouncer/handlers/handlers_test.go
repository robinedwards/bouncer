package handlers_test

import(
    "net/http"
    "net/http/httptest"
    "testing"
	"bouncer/handlers"
	"bouncer/experiment"
)

type MockDB struct {}

func (db MockDB) GetExperiments() []experiment.Experiment {
	tests := make([]experiment.Experiment, 0)

	tests = append(tests, experiment.NewExperiment("test1",
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

	tests = append(tests, experiment.NewExperiment("test2",
		experiment.Alternative{Name: "a", Weight: 1},
		experiment.Alternative{Name: "b", Weight: 1}))

	return tests
}

func TestListExperiments(t *testing.T) {
    mockDb := MockDB{}

    homeHandle := handlers.ListExperiments(mockDb)
    req, _ := http.NewRequest("GET", "", nil)
    w := httptest.NewRecorder()

    homeHandle(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Home page didn't return %v", http.StatusOK)
    }
}

func TestParticipate(t *testing.T) {
    mockDb := MockDB{}

    homeHandle := handlers.Participate(mockDb)
    req, _ := http.NewRequest("GET", "/participate/?uid=1", nil)
    w := httptest.NewRecorder()

    homeHandle(w, req)
    if w.Code != http.StatusOK {
        t.Errorf("Participate page didn't return %v", http.StatusOK)
    }
}

func TestBadParticipate(t *testing.T) {
    mockDb := MockDB{}

    homeHandle := handlers.Participate(mockDb)
    req, _ := http.NewRequest("GET", "/participate/?n=f", nil)
    w := httptest.NewRecorder()

    homeHandle(w, req)
    if w.Code != http.StatusBadRequest {
        t.Errorf("Home page didn't return %v", http.StatusBadRequest)
    }
}

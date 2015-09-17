package handler_test

import(
    "net/http"
    "net/http/httptest"
    "testing"
	"bouncer/handlers"
	"bouncer/abtest"
)

type MockDB struct {}

func (db MockDB) GetABTests() []abtest.ABTest {
	tests := make([]abtest.ABTest, 0)

	tests = append(tests, abtest.NewABTest("test1",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))

	tests = append(tests, abtest.NewABTest("test2",
		abtest.Alternative{Name: "a", Weight: 1},
		abtest.Alternative{Name: "b", Weight: 1}))

	return tests
}

func TestListABTests(t *testing.T) {
    mockDb := MockDB{}

    homeHandle := handlers.ListABTests(mockDb)
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
        t.Errorf("Home page didn't return %v", http.StatusOK)
    }
}

func TestBadParticipate(t *testing.T) {
    mockDb := MockDB{}

    homeHandle := handlers.Participate(mockDb)
    req, _ := http.NewRequest("GET", "/participate/?uid=f", nil)
    w := httptest.NewRecorder()

    homeHandle(w, req)
    if w.Code != http.StatusBadRequest {
        t.Errorf("Home page didn't return %v", http.StatusBadRequest)
    }
}

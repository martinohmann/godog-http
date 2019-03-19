package http_test

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/DATA-DOG/godog"
	godoghttp "github.com/martinohmann/godog-http"
)

func fail(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{"error": http.StatusText(code)})
}

func ok(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fail(w, http.StatusMethodNotAllowed)
		return
	}

	if len(r.Header["X-Auth"]) > 0 && r.Header["X-Auth"][0] == "supersecret" {
		data := make(map[string]interface{})

		if r.Body == nil {
			fail(w, http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			fail(w, http.StatusBadRequest)
			return
		}

		ok(w, data)
		return
	}

	fail(w, http.StatusUnauthorized)
}

func TestMain(m *testing.M) {
	status := godog.RunWithOptions("godog", func(s *godog.Suite) {
		c := godoghttp.NewFeatureContext(http.HandlerFunc(handler))
		c.Register(s)
	}, godog.Options{
		Format: "progress",
		Paths:  []string{"features"},
	})

	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

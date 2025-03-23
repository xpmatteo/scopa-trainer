package handlers

import (
	"net/http"
)

// TestCombinationScenario is an interface for setting up a test scenario for combination captures
type TestCombinationScenario interface {
	SetupCombinationTest()
}

// HandleTestCombination handles the HTTP request to set up a test scenario
type HandleTestCombination struct {
	service TestCombinationScenario
}

// NewHandleTestCombination creates a new test combination handler
func NewHandleTestCombination(service TestCombinationScenario) http.HandlerFunc {
	handler := &HandleTestCombination{
		service: service,
	}
	return handler.ServeHTTP
}

// ServeHTTP handles the HTTP request
func (h *HandleTestCombination) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set up the combination test scenario
	h.service.SetupCombinationTest()

	// Redirect to the main game page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
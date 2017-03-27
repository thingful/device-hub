package pipe

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotConfiguredURIGet404(t *testing.T) {
	t.Parallel()

	router := DefaultRouter()

	req, _ := http.NewRequest("POST", "/abc", nil)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler(router))

	handler.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusNotFound)

}

func TestNonPOSTRequestsGet400(t *testing.T) {
	t.Parallel()

	router := DefaultRouter()

	req, _ := http.NewRequest("GET", "/abc", nil)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler(router))

	handler.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)

}

func TestConfiguredURINoContentGet400(t *testing.T) {
	t.Parallel()

	router := DefaultRouter()

	NewHTTPChannel("/abc", router)

	req, _ := http.NewRequest("POST", "/abc", nil)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler(router))

	handler.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusBadRequest)

}

func TestConfiguredURIContentGet202(t *testing.T) {
	t.Parallel()

	router := DefaultRouter()

	channel := NewHTTPChannel("/abc", router)

	// ensure the 'out' channel is emptied
	go func() { <-channel.Out() }()

	request := "{\"a\" : 1}"

	req, _ := http.NewRequest("POST", "/abc", strings.NewReader(request))

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler(router))

	handler.ServeHTTP(w, req)
	assert.Equal(t, w.Code, http.StatusAccepted)

}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWeatherHandler_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/?cep=01001000", nil)
	w := httptest.NewRecorder()

	WeatherHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestWeatherHandler_InvalidCep(t *testing.T) {
	req := httptest.NewRequest("GET", "/?cep=123", nil)
	w := httptest.NewRecorder()

	WeatherHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", resp.StatusCode)
	}
}

func TestWeatherHandler_CepNotFound(t *testing.T) {
	req := httptest.NewRequest("GET", "/?cep=99999999", nil)
	w := httptest.NewRecorder()

	WeatherHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp.StatusCode)
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestKernelRunHandler_HappyPath(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	body := `{"text":"нужно купить билеты в Шанхай","source":"manual"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/kernel/run", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]any
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := result["EventID"]; !ok {
		t.Errorf("response missing EventID field")
	}
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", w.Header().Get("Content-Type"))
	}
}

func TestKernelRunHandler_EmptyText(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	body := `{"text":"","source":"manual"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/kernel/run", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp errorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if resp.Error == "" {
		t.Errorf("expected non-empty error message")
	}
}

func TestKernelRunHandler_WhitespaceText(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	body := `{"text":"   ","source":"manual"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/kernel/run", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for whitespace-only text, got %d", w.Code)
	}
}

func TestKernelRunHandler_InvalidJSON(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	req := httptest.NewRequest(http.MethodPost, "/v1/kernel/run", bytes.NewBufferString("not-json"))
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestKernelRunHandler_MethodNotAllowed(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	req := httptest.NewRequest(http.MethodGet, "/v1/kernel/run", nil)
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}

func TestKernelRunHandler_MissingSource_DefaultsToAPI(t *testing.T) {
	k := buildKernel()
	h := kernelRunHandler(k)

	body := `{"text":"buy tickets"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/kernel/run", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	h(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

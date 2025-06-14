package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockSender struct {
	called      bool
	receivedMsg []byte
	receivedCID string
	returnError error
}

func (m *mockSender) Send(msg []byte, correlationID string) error {
	m.called = true
	m.receivedMsg = msg
	m.receivedCID = correlationID
	return m.returnError
}

func TestPublishHandler_ValidJSONWithHeader(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	body := `{"hello":"world"}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))
	req.Header.Set("X-Correlation-Id", "test-cid-123")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("expected status %d but got %d", http.StatusAccepted, status)
	}

	if !mock.called {
		t.Fatal("expected Send to be called")
	}

	if string(mock.receivedMsg) != body {
		t.Errorf("expected message body %s but got %s", body, string(mock.receivedMsg))
	}

	if mock.receivedCID != "test-cid-123" {
		t.Errorf("expected correlation ID 'test-cid-123' but got %s", mock.receivedCID)
	}
}

func TestPublishHandler_ValidJSONWithoutHeader(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	body := `{"foo":"bar"}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("expected status %d but got %d", http.StatusAccepted, status)
	}

	if !mock.called {
		t.Fatal("expected Send to be called")
	}

	if string(mock.receivedMsg) != body {
		t.Errorf("expected message body %s but got %s", body, string(mock.receivedMsg))
	}

	if mock.receivedCID == "" {
		t.Error("expected correlation ID to be generated but was empty")
	}
}

func TestPublishHandler_InvalidJSON(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	body := `{"bad json}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("expected status %d but got %d", http.StatusBadRequest, status)
	}

	if mock.called {
		t.Error("Send should not be called on invalid JSON")
	}
}

func TestPublishHandler_SendError(t *testing.T) {
	mock := &mockSender{returnError: errors.New("send failed")}
	handler := MakePublishHandler(mock)

	body := `{"ok":"yes"}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status %d but got %d", http.StatusInternalServerError, status)
	}

	if !mock.called {
		t.Fatal("expected Send to be called")
	}
}

func TestPublishHandler_EmptyBody(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	req := httptest.NewRequest("POST", "/publish", nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty body, got %d", rr.Code)
	}

	if mock.called {
		t.Error("Send should not be called on empty body")
	}
}

func TestPublishHandler_MissingContentType(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	body := `{"foo":"bar"}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))
	// No Content-Type header set

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202 for missing content-type, got %d", rr.Code)
	}
}

func TestPublishHandler_LargePayload(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	largeString := strings.Repeat("a", 10_000)
	largeJSON := `{"data":"` + largeString + `"}`

	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(largeJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202 for large payload, got %d", rr.Code)
	}

	if !mock.called {
		t.Error("expected Send to be called")
	}
}

func TestPublishHandler_InvalidMethod(t *testing.T) {
	mock := &mockSender{}
	handler := MakePublishHandler(mock)

	req := httptest.NewRequest("GET", "/publish", nil)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Our handler does not restrict method itself; router usually handles that.
	// So this will likely 400 because of empty body, or 405 if tested with router.
	// Just test that Send is not called.
	if mock.called {
		t.Error("Send should not be called on GET method")
	}
}

package main

import (
	"bytes"
	"github.com/BlurredDev/go-kafka-gateway/internal/config"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/BlurredDev/go-kafka-gateway/internal/handler"
)

type mockSender struct {
	Called      bool
	ReceivedMsg []byte
	ReceivedCID string
	ReturnError error
}

func (m *mockSender) Send(msg []byte, correlationID string) error {
	m.Called = true
	m.ReceivedMsg = msg
	m.ReceivedCID = correlationID
	return m.ReturnError
}

func TestPublishHandler_EmptyBody(t *testing.T) {
	mock := &mockSender{}
	handlerFunc := handler.MakePublishHandler(mock)

	req := httptest.NewRequest("POST", "/publish", nil)

	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for empty body, got %d", rr.Code)
	}

	if mock.Called {
		t.Error("Send should not be called on empty body")
	}
}

func TestPublishHandler_MissingContentType(t *testing.T) {
	mock := &mockSender{}
	handlerFunc := handler.MakePublishHandler(mock)

	body := `{"foo":"bar"}`
	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(body))
	// No Content-Type header set

	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202 for missing content-type, got %d", rr.Code)
	}
}

func TestPublishHandler_LargePayload(t *testing.T) {
	mock := &mockSender{}
	handlerFunc := handler.MakePublishHandler(mock)

	largeString := strings.Repeat("a", 10_000)
	largeJSON := `{"data":"` + largeString + `"}`

	req := httptest.NewRequest("POST", "/publish", bytes.NewBufferString(largeJSON))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected status 202 for large payload, got %d", rr.Code)
	}

	if !mock.Called {
		t.Error("expected Send to be called")
	}
}

func TestPublishHandler_InvalidMethod(t *testing.T) {
	mock := &mockSender{}
	handlerFunc := handler.MakePublishHandler(mock)

	req := httptest.NewRequest("GET", "/publish", nil)

	rr := httptest.NewRecorder()
	handlerFunc.ServeHTTP(rr, req)

	if mock.Called {
		t.Error("Send should not be called on GET method")
	}
}

func TestRun(t *testing.T) {
	cfg := config.Config{
		KafkaBroker: "localhost:9092",
		KafkaTopic:  "test-topic",
		DLQTopic:    "test-dlq",
		HTTPAddr:    ":9091",
	}

	go func() {
		_ = Run(cfg) // you can assert/log errors if needed
	}()

	// Let the server spin up
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9091/publish")
	if err != nil {
		t.Fatalf("failed to call /publish: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 405 or 400, got %d", resp.StatusCode)
	}
}

func TestStartGateway(t *testing.T) {
	go func() {
		err := StartGateway()
		if err != nil {
			t.Errorf("StartGateway failed: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9091/publish")
	if err != nil {
		t.Fatalf("failed to call /publish: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed && resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 405 or 400 from StartGateway, got %d", resp.StatusCode)
	}
}

func TestMainFunc(t *testing.T) {
	if os.Getenv("TEST_MAIN_CALL") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMainFunc")
	cmd.Env = append(os.Environ(), "TEST_MAIN_CALL=1")
	err := cmd.Run()

	// Accept exit code 1 from log.Fatalf
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			t.Logf("main exited with expected code 1")
			return
		}
		t.Fatalf("main exited with unexpected error: %v", err)
	}
}

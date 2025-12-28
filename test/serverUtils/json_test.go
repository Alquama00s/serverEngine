package serverUtils_test

import (
	"bytes"
	"errors"
	"net/http"
	"testing"

	"github.com/Alquama00s/serverEngine/lib"
	"github.com/Alquama00s/serverEngine/serverUtils"
)

type TestStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestUnmarshal_Success_WithRawBody(t *testing.T) {
	req := &lib.Request{
		RawBody: []byte(`{"name":"test","value":42}`),
	}

	result, err := serverUtils.Unmarshal[TestStruct](req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Name != "test" || result.Value != 42 {
		t.Errorf("Expected Name=test, Value=42, got Name=%s, Value=%d", result.Name, result.Value)
	}

	if req.Body == nil {
		t.Error("Expected req.Body to be set")
	}
}

func TestUnmarshal_Success_ReadFromBody(t *testing.T) {
	body := `{"name":"fromBody","value":123}`
	httpReq, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(body)))
	req := &lib.Request{
		RawRequest: httpReq,
	}

	result, err := serverUtils.Unmarshal[TestStruct](req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result.Name != "fromBody" || result.Value != 123 {
		t.Errorf("Expected Name=fromBody, Value=123, got Name=%s, Value=%d", result.Name, result.Value)
	}

	if req.Body == nil {
		t.Error("Expected req.Body to be set")
	}

	if len(req.RawBody) == 0 {
		t.Error("Expected RawBody to be populated")
	}
}

func TestUnmarshal_ReadBodyError(t *testing.T) {
	// Create a reader that returns an error
	reader := &errorReader{err: errors.New("read error")}
	httpReq, _ := http.NewRequest("POST", "/", reader)
	req := &lib.Request{
		RawRequest: httpReq,
	}

	_, err := serverUtils.Unmarshal[TestStruct](req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err.Error() != "read error" {
		t.Errorf("Expected 'read error', got %v", err)
	}
}

func TestUnmarshal_InvalidJSON(t *testing.T) {
	req := &lib.Request{
		RawBody: []byte(`{"name":"test","value":"invalid"}`), // value should be int
	}

	_, err := serverUtils.Unmarshal[TestStruct](req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

// Helper type for testing read errors
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}

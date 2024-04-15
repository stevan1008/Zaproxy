package implement

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestZAPAdapter_StartScan(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"scan":"12345"}`))
    }))
    defer ts.Close()
    adapter := NewZAPAdapter(ts.URL, "")
    scanID, err := adapter.StartScan([]string{"http://example.com"})
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if scanID != "12345" {
        t.Fatalf("Expected scanID to be '12345', got '%s'", scanID)
    }
}

func TestZAPAdapter_GetScanStatus(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            t.Errorf("Expected 'GET' request, got '%s'", r.Method)
        }

        // if !strings.Contains(r.URL.Path, "/JSON/ascan/view/status/") {
        //     t.Errorf("Unexpected API path: %s", r.URL.Path)
        // }

        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"100"}`)) 
    }))
    defer ts.Close()
    adapter := NewZAPAdapter(ts.URL, "")
    status, err := adapter.GetScanStatus("12345")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if status != "100" {
        t.Errorf("Expected status '100', got '%s'", status)
    }
}

func TestZAPAdapter_GetScanResults(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Query().Get("scanId") != "knownScanID" {
            t.Errorf("Expected to receive 'knownScanID', got '%s'", r.URL.Query().Get("scanId"))
        }

        results := []struct {
            URL    string `json:"url"`
            Alerts []struct {
                Name        string `json:"name"`
                Description string `json:"description"`
                Risk        string `json:"risk"`
            } `json:"alerts"`
        }{
            {
                URL: "http://example.com",
                Alerts: []struct {
                    Name        string `json:"name"`
                    Description string `json:"description"`
                    Risk        string `json:"risk"`
                }{
                    {
                        Name:        "Example Vulnerability",
                        Description: "This is an example vulnerability.",
                        Risk:        "High",
                    },
                },
            },
        }

        respBody, err := json.Marshal(results)
        if err != nil {
            t.Fatalf("Failed to marshal response: %v", err)
        }

        w.WriteHeader(http.StatusOK)
        w.Write(respBody)
    }))
    defer ts.Close()
    adapter := NewZAPAdapter(ts.URL, "")
    results, err := adapter.GetScanResults("knownScanID")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if len(results) != 1 || results[0].URL != "http://example.com" {
        t.Errorf("Unexpected results: %+v", results)
    }
}
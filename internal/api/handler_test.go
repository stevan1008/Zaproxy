package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stevan1008/scanner-service-api-go/pkg/structs"
)

type mockScanner struct {
	mockStartScan   func(targetURLs []string) (string, error)
    mockGetScanStatus func(scanID string) (status string, err error)
	mockGetScanResults func(scanID string) ([]structs.ScanResult, error)
}

func (m *mockScanner) StartScan(targetURLs []string) (scanID string, err error) {
    return m.mockStartScan(targetURLs)
}

func (m *mockScanner) GetScanStatus(scanID string) (status string, err error) {
    return m.mockGetScanStatus(scanID)
}

func (m *mockScanner) GetScanResults(scanID string) ([]structs.ScanResult, error) {
    return m.mockGetScanResults(scanID)
}

func TestStartScan(t *testing.T) {
    scannerMock := &mockScanner{
		mockStartScan: func(targetURLs []string) (string, error) {
			if len(targetURLs) > 0 {
				return "scanID123", nil
			}
			return "", errors.New("error starting scan")
		},
	}

    handler := NewHandler(scannerMock)
    body := strings.NewReader(`{"urls":["http://example.com"]}`)
    req, err := http.NewRequest("POST", "/scan/start", body)

    if err != nil {
        t.Fatal(err)
    }

    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    handler.StartScan(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `{"scan_id":"scanID123"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestGetScanStatus(t *testing.T) {
    scannerMock := &mockScanner{
        mockGetScanStatus: func(scanID string) (string, error) {
            if scanID == "knownScanID" {
                return "Completed", nil
            }
            return "", errors.New("unknown scan ID")
        },
    }

    handler := NewHandler(scannerMock)
    req, err := http.NewRequest("GET", "/scan/status?scanId=knownScanID", nil)
    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()
    handler.GetScanStatus(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

func TestGetScanResults(t *testing.T) {
    scannerMock := &mockScanner{
        mockGetScanResults: func(scanID string) ([]structs.ScanResult, error) {
            if scanID == "knownScanID" {
                results := []structs.ScanResult{
                    {
                        URL: "http://example.com",
                        Alerts: []structs.Alert{
                            {
                                Name:        "Test Vulnerability",
                                Description: "This is a test vulnerability",
                                Risk:        "High",
                            },
                        },
                        Severity: "High",
                    },
                }
                return results, nil
            }
            return nil, errors.New("unknown scan ID")
        },
    }

    handler := NewHandler(scannerMock)
    req, err := http.NewRequest("GET", "/scan/results?scanId=knownScanID", nil)

    if err != nil {
        t.Fatal(err)
    }
    rr := httptest.NewRecorder()

    handler.GetScanResults(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expectedBody := `[{"url":"http://example.com","alerts":[{"name":"Test Vulnerability","description":"This is a test vulnerability","risk":"High"}],"severity":"High"}]`
    if rr.Body.String() != expectedBody {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedBody)
    }
}
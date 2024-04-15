package implement

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/stevan1008/scanner-service-api-go/internal/scanner"
	"github.com/stevan1008/scanner-service-api-go/pkg/structs"
)

var _ scanner.Scanner = &ZAPAdapter{}

type ZAPAdapter struct {
    apiURL string
    apiKey string
}

func NewZAPAdapter(apiURL string, apiKey string) *ZAPAdapter {
    return &ZAPAdapter{
        apiURL: apiURL,
        apiKey: apiKey,
    }
}

// StartScan, start a scan in Zaproxy
func (za *ZAPAdapter) StartScan(targetURLs []string) (string, error) {
	url := fmt.Sprintf("%s/JSON/ascan/action/scan/?apikey=%s&url=%s", za.apiURL, za.apiKey, targetURLs[0])

    request, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request to Zaproxy: %w", err)
    }
    client := &http.Client{Timeout: 10 * time.Second}
	fmt.Printf("Sending request to Zaproxy: %v", request)
    response, err := client.Do(request)
    if err != nil {
        return "", fmt.Errorf("error sending request to Zaproxy: %w", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("zaproxy returned non-200 status code: %s", response.Status)
    }

    var respData map[string]interface{}
    if err := json.NewDecoder(response.Body).Decode(&respData); err != nil {
        return "", fmt.Errorf("error decoding response from Zaproxy: %w", err)
    }

    scanID, ok := respData["scanId"].(string)
    if !ok {
        return "", fmt.Errorf("response from Zaproxy does not contain scanId: %v", respData)
    }

    return scanID, nil
}

// GetScanStatus, get the scan status from Zaproxy
func (za *ZAPAdapter) GetScanStatus(scanID string) (string, error) {
    url := fmt.Sprintf("%s/JSON/ascan/view/status/?scanId=%s&apikey=%s", za.apiURL, scanID, za.apiKey)
    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", fmt.Errorf("error while creating request to get scan status: %w", err)
    }

    resp, err := client.Do(req)
    if err != nil {
        return "", fmt.Errorf("error while sending request to get scan status: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("zaproxy returned non-200 status code: %s", resp.Status)
    }

    var statusResp structs.ZapScanStatusResponse
    if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
        return "", fmt.Errorf("error decoding response from Zaproxy: %w", err)
    }

    return statusResp.Status, nil
}

// GetScanResults, get the scan results from Zaproxy
func (za *ZAPAdapter) GetScanResults(scanID string) ([]structs.ScanResult, error) {
    url := fmt.Sprintf("%s/JSON/ascan/view/scanResults/?scanId=%s&apikey=%s", za.apiURL, scanID, za.apiKey)

    client := &http.Client{Timeout: 10 * time.Second}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("error creating request to get scan results: %w", err)
    }

    resp, err := client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request to get scan results: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("response from Zaproxy is not 200 OK: %s", resp.Status)
    }

    var zapResults structs.ZapScanResultsResponse
    if err := json.NewDecoder(resp.Body).Decode(&zapResults); err != nil {
        return nil, fmt.Errorf("error decoding response from Zaproxy: %w", err)
    }

    // Convert the Zaproxy results to the ScanResult struct
    var scanResults []structs.ScanResult
    for _, result := range zapResults.Results {
        alerts := make([]structs.Alert, len(result.Alerts))
        for i, alert := range result.Alerts {
            alerts[i] = structs.Alert{
                Name:        alert.Name,
                Description: alert.Description,
                Risk:        alert.Risk,
            }
        }
        scanResults = append(scanResults, structs.ScanResult{
            URL:      result.URL,
            Alerts:   alerts,
            Severity: "TBD",
        })
    }

    return scanResults, nil
}
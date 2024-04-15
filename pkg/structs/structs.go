package structs

// Alert represents a security alert found in a scan.
type Alert struct {
    Name        string `json:"name"`
    Description string `json:"description"`
    Risk        string `json:"risk"`
}

// ScanStatus represents the status of a scan.
type ScanStatus struct {
    ID        string `json:"scanId"`
    Status    string `json:"status"`
    Progress  int    `json:"progress"`
}

// ScanResult represents the results of a scan.
type ScanResult struct {
    URL      string   `json:"url"`
    Alerts   []Alert  `json:"alerts"`
    Severity string   `json:"severity"`
}

// ZapScanStatusResponse represents the response from Zaproxy when checking the status of a scan.
type ZapScanStatusResponse struct {
    Status string `json:"status"`
}

// ZapScanResultsResponse represents the response from Zaproxy when getting the results of a scan.
type ZapScanResultsResponse struct {
    Results []struct {
        URL      string   `json:"url"`
        Alerts   []struct {
            Name        string `json:"name"`
            Description string `json:"description"`
            Risk        string `json:"risk"`
        } `json:"alerts"`
    } `json:"results"`
}
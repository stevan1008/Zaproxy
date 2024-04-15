package api

import (
    "encoding/json"
    "net/http"
    "github.com/stevan1008/scanner-service-api-go/internal/scanner"
)

type Handler struct {
    scannerService scanner.Scanner
}

// This is the consumer of the scanner service
func NewHandler(scannerService scanner.Scanner) *Handler {
    return &Handler{
        scannerService: scannerService,
    }
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
    mux.HandleFunc("/scan/start", h.StartScan)
    mux.HandleFunc("/scan/status", h.GetScanStatus)
    mux.HandleFunc("/scan/results", h.GetScanResults)
}

func (h *Handler) StartScan(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }

    var request struct {
        URLs []string `json:"urls"`
    }
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    scanID, err := h.scannerService.StartScan(request.URLs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"scan_id": scanID})
}

func (h *Handler) GetScanStatus(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    scanID := r.URL.Query().Get("scanId")

    if scanID == "" {
        http.Error(w, "Scan ID is required", http.StatusBadRequest)
        return
    }

    status, err := h.scannerService.GetScanStatus(scanID)
    if err != nil {
        http.Error(w, "Failed to get scan status: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(struct {
        Status string `json:"status"`
    }{
        Status: status,
    })
}

func (h *Handler) GetScanResults(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
        return
    }

    scanID := r.URL.Query().Get("scanId")

    if scanID == "" {
        http.Error(w, "Scan ID is required", http.StatusBadRequest)
        return
    }

    results, err := h.scannerService.GetScanResults(scanID)
    if err != nil {
        http.Error(w, "Failed to get scan results: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(results); err != nil {
        http.Error(w, "Failed to encode scan results: "+err.Error(), http.StatusInternalServerError)
    }
}
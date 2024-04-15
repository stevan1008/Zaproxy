package scanner

import (
	"github.com/stevan1008/scanner-service-api-go/pkg/structs"
)

// Scanner, define the methods that a scanner should implement with this interface
type Scanner interface {
    StartScan(targetURLs []string) (scanID string, err error)
    
    GetScanStatus(scanID string) (status string, err error)
    
    GetScanResults(scanID string) (results []structs.ScanResult, err error)
}
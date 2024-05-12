// interfaces/scan_service.go
package interfaces

import "transfigurr/models"

type ScanServiceInterface interface {
	Startup()
	Enqueue(item models.Item)
}

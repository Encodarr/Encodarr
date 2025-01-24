package types

import "transfigurr/internal/interfaces/services"

type Services struct {
	ScanService     services.ScanServiceInterface
	EncodeService   services.EncodeServiceInterface
	MetadataService services.MetadataServiceInterface
}

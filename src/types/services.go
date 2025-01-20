package types

import "transfigurr/interfaces"

type Services struct {
	ScanService     interfaces.ScanServiceInterface
	EncodeService   interfaces.EncodeServiceInterface
	MetadataService interfaces.MetadataServiceInterface
}

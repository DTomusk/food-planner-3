package upload

import "github.com/google/uuid"

type Upload struct {
	Key         string
	URL         string
	Used        bool
	OwnerUserID uuid.UUID
}

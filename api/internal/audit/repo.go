package audit

import (
	"context"
	"foodplanner/internal/db"
)

type repo struct{}

func NewRepo() *repo {
	return &repo{}
}

func (r *repo) Save(ctx context.Context, db db.DBTX, entry *AuditEntry) error {
	return nil
}

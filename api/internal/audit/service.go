package audit

import (
	"context"
	"foodplanner/internal/db"
)

type AuditService struct {
	db   db.DBTX
	repo *repo
}

func NewAuditService(db db.DBTX, repo *repo) *AuditService {
	return &AuditService{
		db:   db,
		repo: repo,
	}
}

func (s *AuditService) Log(ctx context.Context, entry *AuditEntry) error {
	return s.repo.Save(ctx, s.db, entry)
}

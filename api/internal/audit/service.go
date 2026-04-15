package audit

import (
	"context"
	"foodplanner/internal/db"
)

type AuditService struct {
	repo *repo
}

func NewAuditService(repo *repo) *AuditService {
	return &AuditService{
		repo: repo,
	}
}

func (s *AuditService) Log(ctx context.Context, tx db.DBTX, entry *AuditEntry) error {
	return s.repo.Save(ctx, tx, entry)
}

package postgres

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/entity"
	"github.com/Bereke1t2/Memere/memere-backend/internal/domain/repository"
	"github.com/Bereke1t2/Memere/memere-backend/internal/repository/postgres/sqlcgen"
	"github.com/Bereke1t2/Memere/memere-backend/pkg/pagination"
)

// AdminAuditRepo implements repository.AdminAuditRepository.
type AdminAuditRepo struct {
	q *sqlcgen.Queries
}

var _ repository.AdminAuditRepository = (*AdminAuditRepo)(nil)

func NewAdminAuditRepo(pool *pgxpool.Pool) *AdminAuditRepo {
	return &AdminAuditRepo{q: sqlcgen.New(pool)}
}

func (r *AdminAuditRepo) Insert(ctx context.Context, l *entity.AdminAuditLog) error {
	details, err := json.Marshal(l.Details)
	if err != nil {
		log.Printf("admin_audit: marshal details: %v", err)
		details = []byte("{}")
	}
	var pgTargetID pgtype.UUID
	if l.TargetID != nil {
		pgTargetID = toPgUUID(*l.TargetID)
	}
	if err := queriesFor(ctx, r.q).InsertAdminAuditLog(ctx, sqlcgen.InsertAdminAuditLogParams{
		ActorID:    toPgUUID(l.ActorID),
		Action:     l.Action,
		TargetType: l.TargetType,
		TargetID:   pgTargetID,
		Details:    details,
	}); err != nil {
		return err
	}
	return nil
}

func (r *AdminAuditRepo) List(_ context.Context, _ repository.AdminAuditFilter, _ *pagination.Cursor, _ int) ([]*entity.AdminAuditLog, *pagination.Cursor, error) {
	// Stub: full filter logic added in Skill 5.
	return nil, nil, nil
}

package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"log"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"
)

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

func (s *SessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	const op = "SessionRepository.Create"

	sessionModel := toSessionModel(session)
	result := s.DB.WithContext(ctx).Create(sessionModel)
	if err := result.Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, e.ErrSessionTokenHashDuplicate
			}
		}
		return nil, e.WrapDBError(op, err)
	}

	log.Printf("%s: session created successfully", op)
	return toSessionEntity(sessionModel), nil
}

func (s *SessionRepository) GetByID(ctx context.Context, id uint) (*domain.Session, error) {
	const (
		op      = "SessionRepository.GetByID"
		massage = "session found successfully"
	)

	var sessionModel SessionModel
	result := s.DB.WithContext(ctx).First(&sessionModel, id)
	if err := checkGetQueryResult(result, op, massage, e.ErrSessionNotFound); err != nil {
		return nil, err
	}

	return toSessionEntity(&sessionModel), nil
}

func (s *SessionRepository) RevokeSession(ctx context.Context, id string) error {
	const (
		op      = "SessionRepository.RevokeSession"
		message = "session revoked successfully"
	)

	result := s.DB.WithContext(ctx).Model(&SessionModel{}).Where("id = ?", id).Update("is_revoked", true)
	return checkChangeQueryResult(result, op, message, e.ErrSessionNotFound)
}

func (s *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	const (
		op      = "SessionRepository.DeleteSession"
		message = "session deleted successfully"
	)

	result := s.DB.WithContext(ctx).Delete(&SessionModel{}, id)
	return checkChangeQueryResult(result, op, message, e.ErrSessionNotFound)
}

func toSessionModel(s *domain.Session) *SessionModel {
	return &SessionModel{
		Id:               s.Id,
		UserId:           s.UserId,
		RefreshTokenHash: s.RefreshTokenHash,
		IsRevoked:        s.IsRevoked,
		CreatedAt:        s.CreatedAt,
		ExpiresAt:        s.ExpiresAt,
	}
}

func toSessionEntity(s *SessionModel) *domain.Session {
	return &domain.Session{
		Id:               s.Id,
		UserId:           s.UserId,
		RefreshTokenHash: s.RefreshTokenHash,
		IsRevoked:        s.IsRevoked,
		CreatedAt:        s.CreatedAt,
		ExpiresAt:        s.ExpiresAt,
	}
}

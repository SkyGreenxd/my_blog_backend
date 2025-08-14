package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
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
				return nil, e.Wrap(op, e.ErrSessionTokenHashDuplicate)
			}
		}
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(sessionModel), nil
}

func (s *SessionRepository) GetByID(ctx context.Context, id uint) (*domain.Session, error) {
	const op = "SessionRepository.GetByID"
	var sessionModel SessionModel
	result := s.DB.WithContext(ctx).First(&sessionModel, id)
	if err := checkGetQueryResult(result, e.ErrSessionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(&sessionModel), nil
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	const op = "SessionRepository.GetByRefreshToken"
	var sessionModel SessionModel
	result := s.DB.WithContext(ctx).First(&sessionModel, "refresh_token = ?", refreshToken)
	if err := checkGetQueryResult(result, e.ErrSessionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(&sessionModel), nil
}

func (s *SessionRepository) RevokeSession(ctx context.Context, id string) error {
	const op = "SessionRepository.RevokeSession"
	result := s.DB.WithContext(ctx).Model(&SessionModel{}).Where("id = ?", id).Update("is_revoked", true)
	return checkChangeQueryResult(result, op, e.ErrSessionNotFound)
}

func (s *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	const op = "SessionRepository.DeleteSession"
	result := s.DB.WithContext(ctx).Delete(&SessionModel{}, id)
	return checkChangeQueryResult(result, op, e.ErrSessionNotFound)
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

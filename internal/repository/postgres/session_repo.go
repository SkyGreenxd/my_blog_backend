package postgres

import (
	"context"
	"my_blog_backend/internal/domain"
	"my_blog_backend/pkg/e"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRepository struct {
	DB *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		DB: db,
	}
}

// Создание сессии
func (s *SessionRepository) Create(ctx context.Context, session *domain.Session) (*domain.Session, error) {
	const op = "SessionRepository.Create"
	sessionModel := toSessionModel(session)
	result := s.DB.WithContext(ctx).Create(sessionModel)

	if err := postgresDuplicate(result, e.ErrRefreshTokenHashDuplicate); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(sessionModel), nil
}

// Получение сесси с помощью Id сессии
func (s *SessionRepository) GetByID(ctx context.Context, sessionId uint) (*domain.Session, error) {
	const op = "SessionRepository.GetByID"
	var sessionModel SessionModel
	result := s.DB.WithContext(ctx).First(&sessionModel, sessionId)
	if err := checkGetQueryResult(result, e.ErrSessionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(&sessionModel), nil
}

// Получение сесси с помощью Refresh токена для обновления JWT
func (s *SessionRepository) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Session, error) {
	const op = "SessionRepository.GetByRefreshTokenHash"
	var sessionModel SessionModel
	result := s.DB.WithContext(ctx).First(&sessionModel, "refresh_token_hash = ?", refreshTokenHash)
	if err := checkGetQueryResult(result, e.ErrSessionNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toSessionEntity(&sessionModel), nil
}

// Функция аннулирует сессию
// Используется при смене пароля, выхода из аккаунта/всех устройств
func (s *SessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	const op = "SessionRepository.RevokeSession"
	result := s.DB.WithContext(ctx).Model(&SessionModel{}).Where("id = ?", id).Update("is_revoked", true)
	if err := checkChangeQueryResult(result, e.ErrSessionNotFound); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

// Удаление сессии
func (s *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	const op = "SessionRepository.DeleteSession"
	result := s.DB.WithContext(ctx).Delete(&SessionModel{}, id)
	if err := checkChangeQueryResult(result, e.ErrSessionNotFound); err != nil {
		return e.Wrap(op, err)
	}

	return nil
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

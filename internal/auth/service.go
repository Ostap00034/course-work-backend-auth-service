package auth

import (
	"context"
	"errors"
	"time"

	commonpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/common/v1"
	userpbv1 "github.com/Ostap00034/course-work-backend-api-specs/gen/go/user/v1"
	"github.com/Ostap00034/course-work-backend-auth-service/util/jwt"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound          = errors.New("такой пользователь не найден")
	ErrInvalidCredentials    = errors.New("неверные данные")
	ErrTokenInvalidOrExpired = errors.New("токен невалидный или просрочен")
	ErrAccountAlreadyRevoked = errors.New("токен уже удален")
)

type Service interface {
	Login(ctx context.Context, email, password string) (token string, expiresAt time.Time, err error)
	ValidateToken(ctx context.Context, tokenStr string) (*jwt.Claims, *commonpb.UserData, error)
	Revoke(ctx context.Context, tokenStr string) error
}

type service struct {
	repo       Repository
	userClient userpbv1.UserServiceClient
	tokenTTL   time.Duration
}

func NewService(r Repository, uc userpbv1.UserServiceClient, ttl time.Duration) Service {
	return &service{repo: r, userClient: uc, tokenTTL: ttl}
}

func (s *service) Login(ctx context.Context, email, password string) (string, time.Time, error) {
	// Проверяем учетку в UserService
	resp, err := s.userClient.ValidateCredentials(ctx, &userpbv1.ValidateCredentialsRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return "", time.Time{}, ErrInvalidCredentials
	}
	userID, err := uuid.Parse(resp.UserId)
	if err != nil {
		return "", time.Time{}, err
	}
	// Генерируем JWT
	exp := time.Now().Add(s.tokenTTL)
	claims := jwt.NewClaims(userID.String(), resp.Role, exp)
	tok, err := jwt.GenerateToken(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	// Сохраняем в БД
	if _, err := s.repo.CreateToken(ctx, userID, tok, exp); err != nil {
		return "", time.Time{}, err
	}
	return tok, exp, nil
}

func (s *service) ValidateToken(ctx context.Context, tokenStr string) (*jwt.Claims, *commonpbv1.UserData, error) {
	t, err := s.repo.GetToken(ctx, tokenStr)
	if err != nil {
		return nil, nil, err
	}
	if t.Revoked || (t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now())) {
		return nil, nil, ErrTokenInvalidOrExpired
	}

	res, err := s.userClient.GetUserById(ctx, &userpbv1.GetUserByIdRequest{
		UserId: t.UserID.String(),
	})
	if err != nil {
		return nil, nil, err
	}

	claims, err := jwt.ParseToken(tokenStr)
	if err != nil {
		return nil, nil, err
	}
	return claims, res.User, nil
}

func (s *service) Revoke(ctx context.Context, tokenStr string) error {
	t, err := s.repo.GetToken(ctx, tokenStr)
	if err != nil {
		return ErrAccountAlreadyRevoked
	}
	return s.repo.RevokeToken(ctx, t.ID)
}

package auth

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    userpb "github.com/Ostap00034/course-work-backend-user-service/api/user/v1"
    "github.com/Ostap00034/course-work-backend-auth-service/util/jwt"
)

type Service interface {
    Login(ctx context.Context, email, password string) (token string, expiresAt time.Time, err error)
    ValidateToken(ctx context.Context, tokenStr string) (*jwt.Claims, error)
    Revoke(ctx context.Context, tokenStr string) error
}

type service struct {
    repo       Repository
    userClient userpb.UserServiceClient
    tokenTTL   time.Duration
}

func NewService(r Repository, uc userpb.UserServiceClient, ttl time.Duration) Service {
    return &service{repo: r, userClient: uc, tokenTTL: ttl}
}

func (s *service) Login(ctx context.Context, email, password string) (string, time.Time, error) {
    // Проверяем учетку в UserService
    resp, err := s.userClient.ValidateCredentials(ctx, &userpb.ValidateCredentialsRequest{
        Email:    email,
        Password: password,
    })
    if err != nil {
        return "", time.Time{}, errors.New("invalid credentials")
    }
    userID, err := uuid.Parse(resp.UserId)
    if err != nil {
        return "", time.Time{}, err
    }
    // Генерируем JWT
    exp := time.Now().Add(s.tokenTTL)
    claims := jwt.NewClaims(userID.String(), exp)
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

func (s *service) ValidateToken(ctx context.Context, tokenStr string) (*jwt.Claims, error) {
    t, err := s.repo.GetToken(ctx, tokenStr)
    if err != nil {
        return nil, err
    }
    if t.Revoked || (t.ExpiresAt != nil && t.ExpiresAt.Before(time.Now())) {
        return nil, errors.New("token invalid or expired")
    }
    return jwt.ParseToken(tokenStr)
}

func (s *service) Revoke(ctx context.Context, tokenStr string) error {
    t, err := s.repo.GetToken(ctx, tokenStr)
    if err != nil {
        return err
    }
    return s.repo.RevokeToken(ctx, t.ID)
}

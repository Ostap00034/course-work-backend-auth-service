package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Ostap00034/course-work-backend-auth-service/ent"
	"github.com/Ostap00034/course-work-backend-auth-service/ent/token"
	"github.com/google/uuid"
)

type Repository interface {
	CreateToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) (*ent.Token, error)
	GetToken(ctx context.Context, tokenStr string) (*ent.Token, error)
	RevokeToken(ctx context.Context, tokenID uuid.UUID) error
}

type repo struct {
	client *ent.Client
}

func NewRepo(client *ent.Client) Repository {
	return &repo{client: client}
}

func (r *repo) CreateToken(ctx context.Context, userID uuid.UUID, tok string, expiresAt time.Time) (*ent.Token, error) {
	return r.client.Token.
		Create().
		SetUserID(userID).
		SetToken(tok).
		SetIssuedAt(time.Now()).
		SetExpiresAt(expiresAt).
		Save(ctx)
}

func (r *repo) GetToken(ctx context.Context, tokenStr string) (*ent.Token, error) {
	t, err := r.client.Token.
		Query().
		Where(token.TokenEQ(tokenStr)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return nil, errors.New("токен не найден")
	}
	return t, err
}

func (r *repo) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	_, err := r.client.Token.
		UpdateOneID(tokenID).
		SetRevoked(true).
		Save(ctx)
	return err
}

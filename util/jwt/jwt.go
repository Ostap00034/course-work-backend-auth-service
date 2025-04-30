package jwt

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var secret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func NewClaims(userID string, exp time.Time) Claims {
    return Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(exp),
        },
    }
}

func GenerateToken(claims Claims) (string, error) {
    tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tok.SignedString(secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
    parsed, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        if t.Method != jwt.SigningMethodHS256 {
            return nil, errors.New("unexpected signing method")
        }
        return secret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
        return claims, nil
    }
    return nil, errors.New("invalid token")
}

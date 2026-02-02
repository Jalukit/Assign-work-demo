package utils

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"

    "assignwork/src/models"
)

type Claims struct {
    UserID string          `json:"userId"`
    Role   models.UserRole `json:"role"`
    jwt.RegisteredClaims
}

func SignToken(userID primitive.ObjectID, role models.UserRole) (string, error) {
    secret := []byte(os.Getenv("JWT_SECRET"))
    now := time.Now()
    claims := Claims{
        UserID: userID.Hex(),
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}

func ParseToken(tokenStr string) (*Claims, error) {
    secret := []byte(os.Getenv("JWT_SECRET"))
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
        return secret, nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrTokenInvalidClaims
}

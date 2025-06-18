package utils

import (
	"DH52111659-api-quan-ly-suc-khoe/internal/data/models"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	SecretKey string
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{
		SecretKey: secretKey,
	}
}

type TokenClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func (t *TokenService) GenerateTokens(user models.Account) (accessToken, refreshToken string, err error) {
    var wg sync.WaitGroup
    var errAccess, errRefresh error
    accessChan := make(chan string, 1)
    refreshChan := make(chan string, 1)

    wg.Add(2)

    // Goroutine để tạo access token
    go func() {
        defer wg.Done()
        accessClaim := TokenClaims{
            UserID: user.ID.String(),
            Role:   user.Role,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
                IssuedAt:  jwt.NewNumericDate(time.Now()),
                Subject:   user.ID.String(),
            },
        }

        accessJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim)
        token, err := accessJwt.SignedString([]byte(t.SecretKey))
        if err != nil {
            errAccess = err
        } else {
            accessChan <- token
        }
    }()

    // Goroutine để tạo refresh token
    go func() {
        defer wg.Done()
        refreshClaim := TokenClaims{
            UserID: user.ID.String(),
            Role:   user.Role,
            RegisteredClaims: jwt.RegisteredClaims{
                ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
                IssuedAt:  jwt.NewNumericDate(time.Now()),
                Subject:   user.ID.String(),
            },
        }

        refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)
        token, err := refreshJwt.SignedString([]byte(t.SecretKey))
        if err != nil {
            errRefresh = err
        } else {
            refreshChan <- token
        }
    }()

    // Chờ tất cả goroutines hoàn thành
    wg.Wait()
    close(accessChan)
    close(refreshChan)

    // Lấy kết quả từ channels
    if errAccess != nil {
        return "", "", errAccess
    }
    if errRefresh != nil {
        return "", "", errRefresh
    }

    return <-accessChan, <-refreshChan, nil
}

func (s *TokenService) VerifyToken(tokenString string) (*TokenClaims, error) {
	claims := &TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}

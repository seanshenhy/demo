package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v4"
)

type LoginUser struct {
	Email    string
	Username string
	UserID   int
}
type CustomClaims struct {
	LoginUser
	jwt.StandardClaims
}

// GenerateToken create a token string
func GenerateToken(secret []byte, email, username string, userId int) (string, error) {
	claims := CustomClaims{
		LoginUser{
			email,
			username,
			userId,
		},
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(10 * time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

// ParseToken parse custom info
func ParseTokenByCtx(ctx context.Context, secret []byte) (context.Context, error) {
	if tr, ok := transport.FromServerContext(ctx); ok {
		auths := strings.SplitN(tr.RequestHeader().Get("Authorization"), " ", 2)
		if len(auths) != 2 || !strings.EqualFold(auths[0], "Token") {
			return ctx, errors.New(400, "header", "lost jwt token")
		}
		token, err := jwt.ParseWithClaims(auths[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			ctx = context.WithValue(ctx, "loginUser", claims.LoginUser)
			return ctx, nil
		}
		return ctx, err
	}
	return ctx, errors.New(500, "header", "error")
}

// JWTAuth is used for middleware
func JWTAuth(secret []byte) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			ctx, err = ParseTokenByCtx(ctx, secret)
			if err != nil {
				return nil, err
			}
			return handler(ctx, req)
		}
	}
}

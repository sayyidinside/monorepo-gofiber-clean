package helpers

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sayyidinside/gofiber-clean-fresh/services/rbac/domain/entity"
)

func GenerateToken(user *entity.User, expireTime int, secret string, isRefresh bool) (string, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("could not decode token secret: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedSecret)
	if err != nil {
		return "", fmt.Errorf("create: parse token secret key: %w", err)
	}

	claim := make(jwt.MapClaims)
	claim["sub"] = user.ID
	claim["iat"] = time.Now().Unix()
	claim["nbf"] = time.Now().Unix()

	if isRefresh {
		claim["exp"] = time.Now().Add(time.Duration(expireTime) * time.Hour).Unix()
	} else {
		claim["exp"] = time.Now().Add(time.Duration(expireTime) * time.Minute).Unix()

		claim["username"] = user.Username
		claim["email"] = user.Email
		claim["is_admin"] = user.Role.IsAdmin
		claim["validated"] = user.ValidatedAt.Valid
		claim["validated_at"] = user.ValidatedAt.Time.Unix()

		var permissions []string
		for _, permission := range user.Role.Permissions {
			permissions = append(permissions, permission.Name)
		}

		claim["permissions"] = permissions
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claim).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, secret string) (jwt.MapClaims, error) {
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("could not decode token secret: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedSecret)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
			}

			return key, nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claim, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claim, nil
}

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hellokvn/jp-auth-svc/pkg/models"
)

const secretKey string = "secretkey"

type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

type JwtClaims struct {
	jwt.StandardClaims
	Id    int64
	Email string
}

func (j *JwtWrapper) GenerateToken(auth models.Auth) (signedToken string, err error) {
	claims := &JwtClaims{
		Id:    auth.Id,
		Email: auth.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err = token.SignedString([]byte(j.SecretKey))

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (claims *JwtClaims, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JwtClaims)

	if !ok {
		return nil, errors.New("Couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil

}

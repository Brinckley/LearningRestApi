package service

import (
	todo "Todo"
	"Todo/pkg/repository"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	salt       = "sdhkjkl1kjljkljkl21kjl2jjlk1l" // salt for password hash creating
	signingKey = "sdf3kk@qe#o#eq@r4qekjf3#@"
	tokenTTL   = 12 * time.Hour // time to live for token
)

type tokenClaims struct { // creating custom jwt claim, adding own field of userid
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct { // service for authorization (need access to db -> repo object creating)
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService { // DI
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) { // creating user based on received struct
	user.Password = generatePasswordHash(user.Password) // password -> password hash
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New() // using sha1 to create hash of password
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))) // adding salt
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password)) // getting user by name and pass z(check if exists)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ // parameters for implemented struct of jwt.StandardClaims
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey)) // signing token, using secret string
}

func (s *AuthService) ParseToken(accessToken string) (int, error) { // we got the token, parsing info from it
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // checking how token was signed
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims) // getting claims and checking type
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil // returning userId from token
}

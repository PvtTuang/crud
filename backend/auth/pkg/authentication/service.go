package authentication

import (
	"errors"
	"fmt"
	"log"
	"time"

	"database/config"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo UserRepository
	redis    *redis.Client
	jwtKey   []byte
}

func NewAuthService(repo UserRepository, redis *redis.Client, jwtKey []byte) *AuthService {
	return &AuthService{
		userRepo: repo,
		redis:    redis,
		jwtKey:   jwtKey,
	}
}

func (s *AuthService) Register(username, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("username หรือ password ไม่ถูกต้อง")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("username หรือ password ไม่ถูกต้อง")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("auth_token:%d", user.ID)
	err = config.SetToken(s.redis, key, tokenString, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return false, errors.New("token ไม่ถูกต้อง")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("ไม่สามารถอ่าน claims ได้")
	}

	userID := int(claims["user_id"].(float64))
	key := fmt.Sprintf("auth_token:%d", userID)

	valid, err := config.IsTokenValid(s.redis, key, tokenStr)
	if err != nil {
		log.Println("Redis error:", err)
	}
	if !valid {
		return false, errors.New("token ไม่ถูกต้องหรือหมดอายุ")
	}

	return true, nil
}

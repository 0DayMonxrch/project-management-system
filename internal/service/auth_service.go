package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0DayMonxrch/project-management-system/internal/config"
	"github.com/0DayMonxrch/project-management-system/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo domain.UserRepository
	email    domain.EmailService
	cfg      config.JWTConfig
}

func NewAuthService(userRepo domain.UserRepository, email domain.EmailService, cfg config.JWTConfig) domain.AuthService {
	return &authService{userRepo: userRepo, email: email, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, name, email, password string) error {
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return domain.ErrConflict
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	token, err := generateToken()
	if err != nil {
		return err
	}

	user := &domain.User{
		Name:              name,
		Email:             email,
		Password:          string(hash),
		Role:              domain.RoleMember,
		VerificationToken: token,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return s.email.SendVerificationEmail(email, token)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", "", domain.ErrUnauthorized
	}

	if !user.IsEmailVerified {
		return "", "", domain.ErrEmailNotVerified
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", domain.ErrUnauthorized
	}

	accessToken, err := s.generateJWT(user.ID.Hex(), s.cfg.AccessSecret, time.Duration(s.cfg.AccessExpiryMinutes)*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.generateJWT(user.ID.Hex(), s.cfg.RefreshSecret, time.Duration(s.cfg.RefreshExpiryDays)*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) Logout(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	user.RefreshToken = ""
	return s.userRepo.Update(ctx, user)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.FindByVerificationToken(ctx, token)
	if err != nil {
		return domain.ErrTokenInvalid
	}
	user.IsEmailVerified = true
	user.VerificationToken = ""
	return s.userRepo.Update(ctx, user)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := s.parseJWT(refreshToken, s.cfg.RefreshSecret)
	if err != nil {
		return "", domain.ErrTokenInvalid
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return "", domain.ErrTokenInvalid
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", domain.ErrUnauthorized
	}

	if user.RefreshToken != refreshToken {
		return "", domain.ErrTokenInvalid
	}

	return s.generateJWT(userID, s.cfg.AccessSecret, time.Duration(s.cfg.AccessExpiryMinutes)*time.Minute)
}

func (s *authService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return domain.ErrUnauthorized
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	return s.userRepo.Update(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't leak whether email exists
		return nil
	}

	token, err := generateToken()
	if err != nil {
		return err
	}

	user.ResetToken = token
	user.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.email.SendPasswordResetEmail(email, token)
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(ctx, token)
	if err != nil {
		return domain.ErrTokenInvalid
	}

	if time.Now().After(user.ResetTokenExpiry) {
		return domain.ErrTokenExpired
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hash)
	user.ResetToken = ""
	user.ResetTokenExpiry = time.Time{}
	return s.userRepo.Update(ctx, user)
}

func (s *authService) ResendVerificationEmail(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.IsEmailVerified {
		return domain.ErrConflict
	}

	token, err := generateToken()
	if err != nil {
		return err
	}

	user.VerificationToken = token
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return s.email.SendVerificationEmail(user.Email, token)
}

func (s *authService) GetCurrentUser(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// --- helpers ---

func (s *authService) generateJWT(userID, secret string, expiry time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *authService) parseJWT(tokenStr, secret string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}
	return token.Claims, nil
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}
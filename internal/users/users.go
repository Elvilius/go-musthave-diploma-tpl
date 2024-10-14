package users

import (
	"context"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/config"
	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Storer interface {
	GetUserByLogin(ctx context.Context, login string) (models.User, error)
	CreateUser(ctx context.Context, login, passwordHash string) (int, error)
}

type Tokener interface {
	GenerateTokenForUser(userID int) (string, error)
}

type Service struct {
	store Storer
	token Tokener
	cfg   *config.Config
}

func New(store Storer, token Tokener, cfg *config.Config) *Service {
	return &Service{
		store: store,
		token: token,
		cfg:   cfg,
	}
}

func (s *Service) CreateNewUser(ctx context.Context, registerUser models.UserLogin) (models.JWTToken, error) {
	passwordHash, err := s.hashPassword(registerUser.Password)
	if err != nil {
		return "", err
	}
	userID, err := s.store.CreateUser(ctx, registerUser.Login, passwordHash)
	if err != nil {
		return "", err
	}

	token, err := s.token.GenerateTokenForUser(userID)

	return models.JWTToken(token), err
}

func (s *Service) LoginUser(ctx context.Context, login models.UserLogin) (models.JWTToken, error) {
	user, err := s.store.GetUserByLogin(ctx, login.Login)
	if err != nil {
		return "", err
	}

	if ok := s.checkPasswordHash(login.Password, user.PasswordHash); !ok {
		return "", models.ErrUserPasswordNotValid
	}

	token, err := s.token.GenerateTokenForUser(user.ID)

	return models.JWTToken(token), err
}

func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *Service) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

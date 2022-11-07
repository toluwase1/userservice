package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"user-service/config"
	"user-service/db"
	apiError "user-service/errors"
	"user-service/models"
	"user-service/server/jwt"
)

//go:generate mockgen -destination=../mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
// AuthService interface
type AuthService interface {
	IsTokenBlacklisted(token string) bool
	LoginUser(request *models.LoginRequest) (*models.LoginResponse, *apiError.Error)
	SignupUser(request *models.User) (*models.User, *apiError.Error)
}

// authService struct
type authService struct {
	Config   *config.Config
	authRepo db.AuthRepository
}

const Secrete = "aVeryLongSecret"

func (a *authService) IsTokenBlacklisted(token string) bool {
	return a.authRepo.TokenInBlacklist(token)
}

// NewAuthService instantiate an authService
func NewAuthService(authRepo db.AuthRepository, conf *config.Config) AuthService {
	return &authService{
		Config:   conf,
		authRepo: authRepo,
	}
}

func (a *authService) SignupUser(user *models.User) (*models.User, *apiError.Error) {
	err := a.authRepo.IsEmailExist(user.Email)
	if err != nil {
		return nil, apiError.New("email already exist", http.StatusBadRequest)
	}

	err = a.authRepo.IsPhoneExist(user.PhoneNumber)
	if err != nil {
		return nil, apiError.New("phone already exist", http.StatusBadRequest)
	}

	user.HashedPassword, err = GenerateHashPassword(user.Password)
	if err != nil {
		log.Printf("error generating password hash: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	_, err = jwt.GenerateToken(user.Email, int(user.ID), Secrete)
	if err != nil {
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	user.Password = ""
	user.IsEmailActive = true
	user, err = a.authRepo.CreateUser(user)
	if err != nil {
		log.Printf("unable to create user: %v", err.Error())
		return nil, apiError.New("internal server error", http.StatusInternalServerError)
	}
	return user, nil
}

func (a *authService) LoginUser(loginRequest *models.LoginRequest) (*models.LoginResponse, *apiError.Error) {
	foundUser, err := a.authRepo.FindUserByEmail(loginRequest.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apiError.New("invalid email", http.StatusUnprocessableEntity)
		} else {
			log.Printf("error from database: %v", err)
			return nil, apiError.ErrInternalServerError
		}
	}

	if foundUser.IsEmailActive == false {
		return nil, apiError.New("email not verified", http.StatusUnauthorized)
	}

	if err := foundUser.VerifyPassword(loginRequest.Password); err != nil {
		return nil, apiError.ErrInvalidPassword
	}

	accessToken, err := jwt.GenerateToken(foundUser.Email, int(foundUser.ID), "aVeryLongSecret")
	if err != nil {
		log.Printf("error generating token %s", err)
		return nil, apiError.ErrInternalServerError
	}

	return foundUser.LoginUserToDto(accessToken), nil
}

func GenerateHashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

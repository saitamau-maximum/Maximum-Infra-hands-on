package usecase

import (
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/interface/adapter"
	"example.com/webrtc-practice/internal/interface/factory"
)

type UserUsecase struct {
	userRepo      repository.UserRepository
	hasher        adapter.HasherAdapter
	tokenSvc      adapter.TokenServiceAdapter
	userIDFactory factory.UserIDFactory
}

type NewUserUsecaseParams struct {
	UserRepo      repository.UserRepository
	Hasher        adapter.HasherAdapter
	TokenSvc      adapter.TokenServiceAdapter
	UserIDFactory factory.UserIDFactory
}

func NewUserUsecase(p NewUserUsecaseParams) *UserUsecase {
	return &UserUsecase{
		userRepo:      p.UserRepo,
		hasher:        p.Hasher,
		tokenSvc:      p.TokenSvc,
		userIDFactory: p.UserIDFactory,
	}
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse *entity.User

func (u *UserUsecase) SignUp(req SignUpRequest) (SignUpResponse, error) {
	hashedPassword, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	id, err := u.userIDFactory.NewUserID()
	if err != nil {
		return nil, err
	}

	userParams := entity.UserParams{
		ID:         id,
		Name:       req.Name,
		Email:      req.Email,
		PasswdHash: hashedPassword,
		CreatedAt:  time.Now(),
		UpdatedAt:  nil,
	}

	user := entity.NewUser(userParams)

	res, err := u.userRepo.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

type AuthenticateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserResponse string

func (u *UserUsecase) AuthenticateUser(email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	ok, err := u.hasher.ComparePassword(user.GetPasswdHash(), password)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", err
	}

	return u.tokenSvc.GenerateToken(user.GetID())
}

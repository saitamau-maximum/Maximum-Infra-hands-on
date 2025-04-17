package usecase

import (
	"errors"
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

type SignUpResponse struct {
	User *entity.User
}

func (u *UserUsecase) SignUp(req SignUpRequest) (SignUpResponse, error) {
	hashedPassword, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return SignUpResponse{nil}, err
	}

	id, err := u.userIDFactory.NewUserID()
	if err != nil {
		return SignUpResponse{nil}, err
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
		return SignUpResponse{nil}, err
	}

	return SignUpResponse{User: res}, nil
}

type AuthenticateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticateUserResponse struct {
	Token *string `json:"token"`
}
func(res *AuthenticateUserResponse) IsTokenNil() bool {
	return res.Token == nil
}
func (res *AuthenticateUserResponse) GetToken() string {
	if res.Token == nil {
		return ""
	}
	return *res.Token
}

func (u *UserUsecase) AuthenticateUser(req AuthenticateUserRequest) (AuthenticateUserResponse, error) {
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return AuthenticateUserResponse{Token: nil}, err
	}

	ok, err := u.hasher.ComparePassword(user.GetPasswdHash(), req.Password)
	if err != nil {
		return AuthenticateUserResponse{Token: nil}, err
	}

	if !ok {
		return AuthenticateUserResponse{Token: nil}, errors.New("password mismatch")
	}

	res, err := u.tokenSvc.GenerateToken(user.GetID())

	return AuthenticateUserResponse{Token: &res}, err
}

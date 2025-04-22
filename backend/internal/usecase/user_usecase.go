package usecase

import (
	"errors"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/interface/adapter"
	"example.com/webrtc-practice/internal/interface/factory"
)

// UserUseCaseInterface: ユーザーに関するユースケースを管理するインターフェース
type UserUseCaseInterface interface {
	SignUp(req SignUpRequest) (SignUpResponse, error)
	AuthenticateUser(req AuthenticateUserRequest) (AuthenticateUserResponse, error)
	GetUserByID(id entity.UserPublicID) (*entity.User, error)
}

type UserUseCase struct {
	userRepo      repository.UserRepository
	hasher        adapter.HasherAdapter
	tokenSvc      adapter.TokenServiceAdapter
	userIDFactory factory.UserIDFactory
}

type NewUserUseCaseParams struct {
	UserRepo      repository.UserRepository
	Hasher        adapter.HasherAdapter
	TokenSvc      adapter.TokenServiceAdapter
	UserIDFactory factory.UserIDFactory
}

func NewUserUseCase(p NewUserUseCaseParams) *UserUseCase {
	return &UserUseCase{
		userRepo:      p.UserRepo,
		hasher:        p.Hasher,
		tokenSvc:      p.TokenSvc,
		userIDFactory: p.UserIDFactory,
	}
}

// SignUpRequest構造体: サインアップリクエスト
type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse構造体: サインアップレスポンス
type SignUpResponse struct {
	User *entity.User
}

// SignUp ユーザー登録
func (u *UserUseCase) SignUp(req SignUpRequest) (SignUpResponse, error) {
	hashedPassword, err := u.hasher.HashPassword(req.Password)
	if err != nil {
		return SignUpResponse{nil}, err
	}

	publicID, err := u.userIDFactory.NewUserPublicID()
	if err != nil {
		return SignUpResponse{nil}, err
	}

	userParams := entity.UserParams{
		ID:         -1, // IDはDBに保存後に更新されるため、-1を指定
		PublicID:   publicID,
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

// AuthenticateUserRequest構造体: 認証リクエスト
type AuthenticateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthenticateUserResponse構造体: 認証レスポンス
type AuthenticateUserResponse struct {
	token *string
}

// IsTokenNil: トークンがnilかどうかを判定
func (res *AuthenticateUserResponse) IsTokenNil() bool {
	return res.token == nil
}

// GetToken: トークンを取得（nilの場合は空文字を返す）
func (res *AuthenticateUserResponse) GetToken() string {
	if res.token == nil {
		return ""
	}
	return *res.token
}

// 外部でのテストのためのセッター
func (res *AuthenticateUserResponse) SetToken(token string) {
	res.token = &token
}

// AuthenticateUser ユーザー認証
func (u *UserUseCase) AuthenticateUser(req AuthenticateUserRequest) (AuthenticateUserResponse, error) {
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	ok, err := u.hasher.ComparePassword(user.GetPasswdHash(), req.Password)
	if err != nil {
		return AuthenticateUserResponse{token: nil}, err
	}

	if !ok {
		return AuthenticateUserResponse{token: nil}, errors.New("password mismatch")
	}

	res, err := u.tokenSvc.GenerateToken(user.GetID())

	return AuthenticateUserResponse{token: &res}, err
}

func (u *UserUseCase) GetUserByID(id entity.UserPublicID) (*entity.User, error) {
	userID, err := u.userRepo.GetIDByPublicID(id)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

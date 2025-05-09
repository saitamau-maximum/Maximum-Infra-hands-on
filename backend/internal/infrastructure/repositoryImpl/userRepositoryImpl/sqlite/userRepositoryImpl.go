package sqliteuserrepoimpl

import (
	"context"
	"errors"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/model"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

type NewUserRepositoryImplParams struct {
	DB *sqlx.DB
}

func (p *NewUserRepositoryImplParams) Validate() error {
	if p.DB == nil {
		return errors.New("db cannot be nil")
	}
	return nil
}

func NewUserRepositoryImpl(params *NewUserRepositoryImplParams) repository.UserRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}
	return &UserRepositoryImpl{
		db: params.DB,
	}
}

func (r *UserRepositoryImpl) SaveUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		string(user.GetID()),
		user.GetName(),
		user.GetEmail(),
		user.GetPasswdHash(),
		user.GetCreatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetUserByID(ctx, user.GetID())
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id entity.UserID) (*entity.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be 0")
	}

	row := r.db.QueryRowxContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?`, id)

	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	row := r.db.QueryRowxContext(ctx, `
		SELECT id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?`, email)

	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entity.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	// UserID -> UUID
	id := user.GetID()
	userUUID, err := id.UserID2UUID()
	if err != nil {
		return err
	}
	// 現在時刻取得
	time := time.Now()
	var imagePath *string
	path, ok := user.GetImagePath()
	if !ok {
		imagePath = nil
	} else {
		imagePath = &path
	}

	userModel := model.UserModel{
		ID:           userUUID,
		Name:         user.GetName(),
		Email:        user.GetEmail(),
		PasswordHash: user.GetPasswdHash(),
		ImagePath:    imagePath,
		CreatedAt:    user.GetCreatedAt(),
		UpdatedAt:    &time,
	}

	_, err = r.db.NamedExecContext(ctx, `
		UPDATE users
		SET name = :name,
			email = :email,
			password_hash = :password_hash,
			image_path = :image_path,
			updated_at = :updated_at
			WHERE id = :id`, &userModel)

	if err != nil {
		return err
	}
	return nil
}

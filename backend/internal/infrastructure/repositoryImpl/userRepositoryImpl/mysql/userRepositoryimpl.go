package mysqluserrepoimpl

import (
	"errors"

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

func (r *UserRepositoryImpl) SaveUser(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, errors.New("user cannot be nil")
	}

	_, err := r.db.Exec(`
		INSERT INTO users (public_id, name, email, password_hash, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		string(user.GetPublicID()),
		user.GetName(),
		user.GetEmail(),
		user.GetPasswdHash(),
		user.GetCreatedAt(),
	)
	if err != nil {
		return nil, err
	}

	return r.GetUserByPublicID(user.GetPublicID())
}

func (r *UserRepositoryImpl) GetUserByID(id entity.UserID) (*entity.User, error) {
	if id == 0 {
		return nil, errors.New("id cannot be 0")
	}

	row := r.db.QueryRowx(`
		SELECT id, public_id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = ?`, id)

	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	row := r.db.QueryRowx(`
		SELECT id, public_id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?`, email)

	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *UserRepositoryImpl) GetUserByPublicID(publicID entity.UserPublicID) (*entity.User, error) {
	if publicID == "" {
		return nil, errors.New("publicID cannot be empty")
	}

	row := r.db.QueryRowx(`
		SELECT id, public_id, name, email, password_hash, created_at, updated_at
		FROM users
		WHERE public_id = ?`, string(publicID))

	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}

	return userModel.ToEntity(), nil
}

func (r *UserRepositoryImpl) GetIDByPublicID(publicID entity.UserPublicID) (entity.UserID, error) {
	if publicID == "" {
		return 0, errors.New("publicID cannot be empty")
	}

	row := r.db.QueryRowx(`SELECT id FROM users WHERE public_id = ?`, string(publicID))

	var id int64
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return entity.UserID(id), nil
}

func (r *UserRepositoryImpl) GetPublicIDByID(id entity.UserID) (entity.UserPublicID, error) {
	if id == 0 {
		return "", errors.New("id cannot be 0")
	}

	row := r.db.QueryRowx(`SELECT public_id FROM users WHERE id = ?`, id)

	var publicID string
	if err := row.Scan(&publicID); err != nil {
		return "", err
	}

	return entity.UserPublicID(publicID), nil
}

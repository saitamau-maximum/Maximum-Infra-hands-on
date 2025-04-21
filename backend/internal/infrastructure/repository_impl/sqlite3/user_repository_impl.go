package sqlite3

import (
	"errors"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryImpl struct {
	db *sqlx.DB
}

type NewUserRepositoryParams struct {
	DB *sqlx.DB
}

func (p *NewUserRepositoryParams) Validate() error {
	if p.DB == nil {
		return errors.New("DB is required")
	}
	return nil
}

func NewUserRepository(params NewUserRepositoryParams) repository.UserRepository {
	if err := params.Validate(); err != nil {
		panic(err)
	}

	return &UserRepositoryImpl{
		db: params.DB,
	}
}

type UserModel struct {
	ID         string     `db:"id"`
	Name       string     `db:"name"`
	Email      string     `db:"email"`
	PasswdHash string     `db:"passwd_hash"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (r *UserRepositoryImpl) SaveUser(user *entity.User) (*entity.User, error) {
	// Insert user into the database
	query := `INSERT INTO users (name, email, passwd_hash) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, user.GetName(), user.GetEmail(), user.GetPasswdHash())
	if err != nil {
		return nil, err
	}

	// Retrieve the inserted user
	query = `SELECT id, name, email, passwd_hash, created_at FROM users WHERE email = ?`
	var userModel UserModel
	err = r.db.Get(&userModel, query, user.GetEmail())
	if err != nil {
		return nil, err
	}

	res := entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswdHash,
		CreatedAt:  userModel.CreatedAt,
		UpdatedAt:  userModel.UpdatedAt,
	})

	return res, nil
}

func (r *UserRepositoryImpl) GetUserByID(id entity.UserID) (*entity.User, error) {
	// Retrieve user by ID from the database
	query := `SELECT id, name, email, passwd_hash, created_at, updated_at FROM users WHERE id = ?`
	var userModel UserModel
	err := r.db.Get(&userModel, query, id)
	if err != nil {
		return nil, err
	}

	res := entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswdHash,
		CreatedAt:  userModel.CreatedAt,
		UpdatedAt:  userModel.UpdatedAt,
	})

	return res, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	// Retrieve user by email from the database
	query := `SELECT id, name, email, passwd_hash, created_at, updated_at FROM users WHERE email = ?`
	var userModel UserModel
	err := r.db.Get(&userModel, query, email)
	if err != nil {
		return nil, err
	}

	res := entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswdHash,
		CreatedAt:  userModel.CreatedAt,
		UpdatedAt:  userModel.UpdatedAt,
	})

	return res, nil
}

package sqliteuserrepoimpl

import (
	"errors"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/domain/repository"
	"example.com/webrtc-practice/internal/infrastructure/repository_impl/model"
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

	_, err := r.db.Exec(`INSERT INTO users (public_id, name, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?)`,
		string(user.GetPublicID()),
		user.GetName(),
		user.GetEmail(),
		user.GetPasswdHash(),
		user.GetCreatedAt().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, err
	}

	// publicIDを使ってユーザーを取得
	row := r.db.QueryRowx(`SELECT id, public_id, name, email, password_hash, created_at, updated_at FROM users WHERE public_id = ?`, string(user.GetPublicID()))
	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}
	// model.UserModelからUserを生成
	// userModel.CreatedAtをtime.Timeに変換
	createdAt, err := time.Parse("2006-01-02 15:04:05", userModel.CreatedAt)
	if err != nil {
		return nil, err
	}
	// userModel.UpdatedAtを*time.Timeに変換
	// nilの場合はnilを代入
	var updatedAt *time.Time
	if userModel.UpdatedAt != nil {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", *userModel.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt = &parsedTime
	}
	user = entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		PublicID:   entity.UserPublicID(userModel.PublicID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswordHash,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	})

	return user, nil
}

func (r *UserRepositoryImpl) GetUserByID(id entity.UserID) (*entity.User, error) {
	if id == 0 {
		return nil, errors.New("id cannot be 0")
	}

	row := r.db.QueryRowx(`SELECT id, public_id, name, email, password_hash, created_at, updated_at FROM users WHERE id = ?`, id)
	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}
	// model.UserModelからUserを生成
	// userModel.CreatedAtをtime.Timeに変換
	createdAt, err := time.Parse("2006-01-02 15:04:05", userModel.CreatedAt)
	if err != nil {
		return nil, err
	}
	// userModel.UpdatedAtを*time.Timeに変換
	// nilの場合はnilを代入
	var updatedAt *time.Time
	if userModel.UpdatedAt != nil {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", *userModel.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt = &parsedTime
	}
	user := entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		PublicID:   entity.UserPublicID(userModel.PublicID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswordHash,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	})

	return user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	row := r.db.QueryRowx(`SELECT id, public_id, name, email, password_hash, created_at, updated_at FROM users WHERE email = ?`, email)
	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return nil, err
	}
	// model.UserModelからUserを生成
	// userModel.CreatedAtをtime.Timeに変換
	createdAt, err := time.Parse("2006-01-02 15:04:05", userModel.CreatedAt)
	if err != nil {
		return nil, err
	}
	// userModel.UpdatedAtを*time.Timeに変換
	// nilの場合はnilを代入
	var updatedAt *time.Time
	if userModel.UpdatedAt != nil {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", *userModel.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt = &parsedTime
	}
	user := entity.NewUser(entity.UserParams{
		ID:         entity.UserID(userModel.ID),
		PublicID:   entity.UserPublicID(userModel.PublicID),
		Name:       userModel.Name,
		Email:      userModel.Email,
		PasswdHash: userModel.PasswordHash,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	})

	return user, nil
}

func (r *UserRepositoryImpl) GetIDByPublicID(publicID entity.UserPublicID) (entity.UserID, error) {
	if publicID == "" {
		return 0, errors.New("publicID cannot be empty")
	}

	row := r.db.QueryRowx(`SELECT id FROM users WHERE public_id = ?`, string(publicID))
	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return 0, err
	}

	return entity.UserID(userModel.ID), nil
}

func (r *UserRepositoryImpl) GetPublicIDByID(id entity.UserID) (entity.UserPublicID, error) {
	if id == 0 {
		return "", errors.New("id cannot be 0")
	}

	row := r.db.QueryRowx(`SELECT public_id FROM users WHERE id = ?`, id)
	var userModel model.UserModel
	if err := row.StructScan(&userModel); err != nil {
		return "", err
	}

	return entity.UserPublicID(userModel.PublicID), nil
}

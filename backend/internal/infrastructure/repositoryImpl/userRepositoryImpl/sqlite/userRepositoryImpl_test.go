package sqliteuserrepoimpl_test

import (
	"context"
	"testing"
	"time"

	"example.com/infrahandson/internal/domain/entity"
	sqlitegatewayimpl "example.com/infrahandson/internal/infrastructure/gatewayImpl/db/sqlite"
	sqliteuserrepoimpl "example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// SQLiteインメモリDBをセットアップ
	initializer := sqlitegatewayimpl.NewSQLiteInitializer(&sqlitegatewayimpl.NewSQLiteInitializerParams{
		Path:           ":memory:",
		MigrationsPath: "file://migrations",
	})
	db, err := initializer.Init()
	require.NoError(t, err)

	// テスト用スキーマを作成
	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		public_id TEXT NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME
	);
	CREATE TABLE rooms (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  public_id TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL
	);
	CREATE TABLE room_members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id INTEGER NOT NULL,
		user_id TEXT NOT NULL,
		FOREIGN KEY (room_id) REFERENCES rooms(id)
	);`)
	require.NoError(t, err)

	return db
}

func TestSaveUser(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		ID:   "test-public-id",
		Name:       "John Doe",
		Email:      "johndoe@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(context.Background(), user)
	require.NoError(t, err)

	// 保存されたユーザーを取得して確認
	fetchedUser, err := userRepo.GetUserByID(context.Background(), savedUser.GetID())
	require.NoError(t, err)

	// fetchedUserとsavedUserが一致することを確認
	assert.Equal(t, savedUser.GetID(), fetchedUser.GetID())
	assert.Equal(t, savedUser.GetName(), fetchedUser.GetName())
	assert.Equal(t, savedUser.GetEmail(), fetchedUser.GetEmail())
}

func TestGetUserByEmail(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		ID:   "test-public-id",
		Name:       "Jane Doe",
		Email:      "janedoe@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(context.Background(), user)
	require.NoError(t, err)

	// 保存したユーザーをEmailで取得
	fetchedUser, err := userRepo.GetUserByEmail(context.Background(), savedUser.GetEmail())
	require.NoError(t, err)

	// fetchedUserとsavedUserが一致することを確認
	assert.Equal(t, savedUser.GetID(), fetchedUser.GetID())
	assert.Equal(t, savedUser.GetName(), fetchedUser.GetName())
}

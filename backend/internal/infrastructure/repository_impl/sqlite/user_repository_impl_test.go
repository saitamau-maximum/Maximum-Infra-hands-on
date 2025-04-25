package sqlite3_test

import (
	"testing"
	"time"

	"example.com/webrtc-practice/internal/domain/entity"
	"example.com/webrtc-practice/internal/infrastructure/gateway_impl"
	sqlite3 "example.com/webrtc-practice/internal/infrastructure/repository_impl/sqlite"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	// SQLiteインメモリDBをセットアップ
	initializer := gateway_impl.NewSQLiteInitializer(":memory:")
	db, err := initializer.Init()
	require.NoError(t, err)

	// テスト用スキーマを作成
	_, err = db.Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		public_id TEXT NOT NULL,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT
	);`)
	require.NoError(t, err)

	return db
}

func TestSaveUser(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		ID: entity.UserID(-1),
		PublicID:   "test-public-id",
		Name:       "John Doe",
		Email:      "johndoe@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// 保存されたユーザーを取得して確認
	fetchedUser, err := userRepo.GetUserByID(savedUser.GetID())
	require.NoError(t, err)

	// fetchedUserとsavedUserが一致することを確認
	assert.Equal(t, savedUser.GetPublicID(), fetchedUser.GetPublicID())
	assert.Equal(t, savedUser.GetName(), fetchedUser.GetName())
	assert.Equal(t, savedUser.GetEmail(), fetchedUser.GetEmail())
}

func TestGetUserByEmail(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		PublicID:   "test-public-id",
		Name:       "Jane Doe",
		Email:      "janedoe@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// 保存したユーザーをEmailで取得
	fetchedUser, err := userRepo.GetUserByEmail(savedUser.GetEmail())
	require.NoError(t, err)

	// fetchedUserとsavedUserが一致することを確認
	assert.Equal(t, savedUser.GetPublicID(), fetchedUser.GetPublicID())
	assert.Equal(t, savedUser.GetName(), fetchedUser.GetName())
}

func TestGetIDByPublicID(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		PublicID:   "test-public-id",
		Name:       "Alice",
		Email:      "alice@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// 保存したユーザーのPublicIDを使ってIDを取得
	userID, err := userRepo.GetIDByPublicID(savedUser.GetPublicID())
	require.NoError(t, err)

	// 取得したIDが保存したユーザーのIDと一致することを確認
	assert.Equal(t, savedUser.GetID(), userID)
}

func TestGetPublicIDByID(t *testing.T) {
	// テスト用DBをセットアップ
	db := setupTestDB(t)
	defer db.Close()

	// UserRepositoryをインスタンス化
	userRepo := sqlite3.NewUserRepositoryImpl(&sqlite3.NewUserRepositoryImplParams{DB: db})

	// テスト用ユーザーを作成
	user := entity.NewUser(entity.UserParams{
		PublicID:   "test-public-id",
		Name:       "Bob",
		Email:      "bob@example.com",
		PasswdHash: "hashed-password",
		CreatedAt:  time.Now(),
	})

	// ユーザーを保存
	savedUser, err := userRepo.SaveUser(user)
	require.NoError(t, err)

	// 保存したユーザーのIDを使ってPublicIDを取得
	publicID, err := userRepo.GetPublicIDByID(savedUser.GetID())
	require.NoError(t, err)

	// 取得したPublicIDが保存したユーザーのPublicIDと一致することを確認
	assert.Equal(t, savedUser.GetPublicID(), publicID)
}

package di

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/mysqlmsgrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlitemsgrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/mysqlroomrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/roomRepositoryImpl/sqliteroomrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/mysqluserrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/userRepositoryImpl/sqliteuserrepo"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/websocketClientRepositoryImpl/memwsclientrepo"
	"github.com/jmoiron/sqlx"
)

// DBType はデータベースの種類を表す型です。
// repository のDI組み立てで、引数として受け取る
// 安全性のために用意している
type DBType string

const (
	DBTypeMySQL  DBType = "mysql"
	DBTypeSQLite DBType = "sqlite"
)

func RepoistoryInitialize(
	dbType DBType,
	db *sqlx.DB,
) repository.Repository {
	var userRepository repository.UserRepository
	var roomRepository repository.RoomRepository
	var msgRepository repository.MessageRepository

	// Repositoryの初期化
	switch dbType {
	case DBTypeMySQL:
		userRepository = mysqluserrepo.NewUserRepositoryImpl(&mysqluserrepo.NewUserRepositoryImplParams{DB: db})
		roomRepository = mysqlroomrepo.NewRoomRepositoryImpl(&mysqlroomrepo.NewRoomRepositoryImplParams{DB: db})
		msgRepository = mysqlmsgrepo.NewMessageRepositoryImpl(&mysqlmsgrepo.NewMessageRepositoryImplParams{DB: db})
	case DBTypeSQLite:
		userRepository = sqliteuserrepo.NewUserRepositoryImpl(&sqliteuserrepo.NewUserRepositoryImplParams{DB: db})
		roomRepository = sqliteroomrepo.NewRoomRepositoryImpl(&sqliteroomrepo.NewRoomRepositoryImplParams{DB: db})
		msgRepository = sqlitemsgrepo.NewMessageRepositoryImpl(&sqlitemsgrepo.NewMessageRepositoryImplParams{DB: db})
	}

	wsClientRepository := memwsclientrepo.NewInMemoryWebsocketClientRepository(memwsclientrepo.NewInMemoryWebsocketClientRepositoryParams{})

	return repository.Repository{
		UserRepository:     userRepository,
		RoomRepository:     roomRepository,
		MessageRepository:  msgRepository,
		WsClientRepository: wsClientRepository,
	}
}

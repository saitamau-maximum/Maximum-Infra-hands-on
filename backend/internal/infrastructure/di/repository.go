package di

import (
	"example.com/infrahandson/internal/domain/repository"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/mysqlmsgrepoimpl"
	"example.com/infrahandson/internal/infrastructure/repositoryImpl/messageRepositoryImpl/sqlitemsgrepoimpl"
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
	if dbType == DBTypeMySQL {
		userRepository = mysql.NewUserRepositoryImpl(&mysqluserrepoimpl.NewUserRepositoryImplParams{DB: db})
		roomRepository = mysqlroomrepoimpl.NewRoomRepositoryImpl(&mysqlroomrepoimpl.NewRoomRepositoryImplParams{DB: db})
		msgRepository = mysqlmsgrepoimpl.NewMessageRepositoryImpl(&mysqlmsgrepoimpl.NewMessageRepositoryImplParams{DB: db})
	} else if dbType == DBTypeSQLite {
		userRepository = sqliteuserrepoimpl.NewUserRepositoryImpl(&sqliteuserrepoimpl.NewUserRepositoryImplParams{DB: db})
		roomRepository = sqliteroomrepoimpl.NewRoomRepositoryImpl(&sqliteroomrepoimpl.NewRoomRepositoryImplParams{DB: db})
		msgRepository = sqlitemsgrepoimpl.NewMessageRepositoryImpl(&sqlitemsgrepoimpl.NewMessageRepositoryImplParams{DB: db})
	}

	wsClientRepository := inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepository(inmemorywsclientrepoimpl.NewInMemoryWebsocketClientRepositoryParams{})

	return repository.Repository{
		UserRepository:     userRepository,
		RoomRepository:     roomRepository,
		MessageRepository:  msgRepository,
		WsClientRepository: wsClientRepository,
	}
}

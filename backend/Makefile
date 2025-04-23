.PHONY: test-usecase test-user-usecase test-room-usecase test-websocket-usecase

test: test-usecase test-handler

test-usecase: test-user-usecase test-room-usecase test-websocket-usecase test-message-usecase

test-user-usecase:
	go test -v ./internal/usecase/user_usecase_test.go

test-room-usecase:
	go test -v ./internal/usecase/room_usecase_test.go

test-websocket-usecase:
	go test -v ./internal/usecase/websocket_usecase_test.go

test-message-usecase:
	go test -v ./internal/usecase/message_usecase_test.go

test-handler: test-user-handler test-room-handler test-websocket-handler

test-user-handler:
	go test -v ./internal/interface/handler/user_handler_test.go

test-room-handler:
	go test -v ./internal/interface/handler/room_handler_test.go

test-websocket-handler:
	go test -v ./internal/interface/handler/websocket_handler_test.go
	
help:
	@echo "Usage:"
	@echo "  make test           		# Run all tests"
	@echo "  make test-user-usecase 	# Run only UserUseCase tests"
	@echo "  make test-signup    		# Run only SignUp tests"
	@echo "  make test-auth      		# Run only AuthenticateUser tests"
	@echo "  make test-room-usecase	# Run only RoomUseCase tests"

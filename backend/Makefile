.PHONY: test-usecase test-user-usecase test-room-usecase test-websocket-usecase

test-usecase: test-user-usecase test-room-usecase test-websocket-usecase

test-user-usecase:
	go test -v /home/nakamuraitsuki/program/Maximum/Infra/Maximum-Infra-hands-on/backend/internal/usecase/user_usecase_test.go

test-room-usecase:
	go test -v /home/nakamuraitsuki/program/Maximum/Infra/Maximum-Infra-hands-on/backend/internal/usecase/room_usecase_test.go

test-websocket-usecase:
	go test -v /home/nakamuraitsuki/program/Maximum/Infra/Maximum-Infra-hands-on/backend/internal/usecase/websocket_usecase_test.go

test-handler: test-user-handler test-room-handler

test-user-handler:
	go test -v /home/nakamuraitsuki/program/Maximum/Infra/Maximum-Infra-hands-on/backend/internal/interface/handler/user_handler_test.go

test-room-handler:
	go test -v /home/nakamuraitsuki/program/Maximum/Infra/Maximum-Infra-hands-on/backend/internal/interface/handler/room_handler_test.go

help:
	@echo "Usage:"
	@echo "  make test           		# Run all tests"
	@echo "  make test-user-usecase 	# Run only UserUseCase tests"
	@echo "  make test-signup    		# Run only SignUp tests"
	@echo "  make test-auth      		# Run only AuthenticateUser tests"
	@echo "  make test-room-usecase	# Run only RoomUseCase tests"

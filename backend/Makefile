.PHONY: test-user-usecase test-room-usecase

test-user_usecase:
	go test -v ./internal/usecase/user_usecase_test.go

test-room-usecase:
	go test -v ./internal/usecase/room_usecase_test.go

help:
	@echo "Usage:"
	@echo "  make test           		# Run all tests"
	@echo "  make test-user-usecase 	# Run only UserUseCase tests"
	@echo "  make test-signup    		# Run only SignUp tests"
	@echo "  make test-auth      		# Run only AuthenticateUser tests"
	@echo "  make test-room-usecase	# Run only RoomUseCase tests"

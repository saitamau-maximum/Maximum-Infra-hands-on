.PHONY: test test-user-usecase test-signup test-auth

test:
	go test -v ./...

test-user_usecase:
	go test -v ./internal/usecase/user_usecase_test.go

test-signup:
	go test -v -run '^TestSignUp$$' ./internal/usecase/user_usecase_test.go

test-auth:
	go test -v -run '^TestAuthenticateUser$$' ./internal/usecase/user_usecase_test.go

help:
	@echo "Usage:"
	@echo "  make test           		# Run all tests"
	@echo "  make test-user-usecase 	# Run only UserUseCase tests"
	@echo "  make test-signup    		# Run only SignUp tests"
	@echo "  make test-auth      		# Run only AuthenticateUser tests"

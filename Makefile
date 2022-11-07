up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/auth_mock.go -package=mocks github.com/decagonhq/meddle-api/services AuthService
	 mockgen -destination=mocks/auth_repo_mock.go -package=mocks github.com/decagonhq/meddle-api/db AuthRepository

test: generate-mock
	 MEDDLE_ENV=test go test ./...

-include .env

.PHONY: test lint lint-fix \
	proto-format proto-lint proto-gen \
	mock-gen \
	migrate-up migrate-up-one migrate-down migrate-reset migrate-create migrate-status

test:
	mkdir -p tmp
	go test -cover ./internal/... -coverprofile=./tmp/cover.out
	go tool cover -func=./tmp/cover.out | tail -n 1
	go tool cover -html=./tmp/cover.out -o ./tmp/cover.html
	open ./tmp/cover.html

lint:
	go tool golangci-lint run ./...

lint-fix:
	go tool golangci-lint run --fix ./...

proto-format:
	buf format -w

proto-lint:
	buf lint

proto-gen:
	buf generate

mock-gen:
	go tool mockgen -source=internal/domain/customer/customer.go -destination=mocks/domain/customer/mock_customer.go -package=mocks
	go tool mockgen -source=internal/domain/customer/repository.go -destination=mocks/domain/customer/mock_repository.go -package=mocks
	go tool mockgen -source=internal/application/customer/usecase.go -destination=mocks/application/customer/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/domain/measurement/measurement.go -destination=mocks/domain/measurement/mock_measurement.go -package=mocks
	go tool mockgen -source=internal/domain/measurement/measurement_entry.go -destination=mocks/domain/measurement/mock_measurement_entry.go -package=mocks
	go tool mockgen -source=internal/domain/measurement/repository.go -destination=mocks/domain/measurement/mock_repository.go -package=mocks
	go tool mockgen -source=internal/application/measurement/usecase.go -destination=mocks/application/measurement/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/domain/measurementitem/measurement_item.go -destination=mocks/domain/measurementitem/mock_measurement_item.go -package=mocks
	go tool mockgen -source=internal/domain/measurementitem/repository.go -destination=mocks/domain/measurementitem/mock_repository.go -package=mocks
	go tool mockgen -source=internal/application/measurementitem/usecase.go -destination=mocks/application/measurementitem/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/domain/organization/organization.go -destination=mocks/domain/organization/mock_organization.go -package=mocks
	go tool mockgen -source=internal/domain/organization/repository.go -destination=mocks/domain/organization/mock_repository.go -package=mocks
	go tool mockgen -source=internal/application/organization/usecase.go -destination=mocks/application/organization/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/domain/standard/age_group_standard.go -destination=mocks/domain/standard/mock_age_group_standard.go -package=mocks
	go tool mockgen -source=internal/domain/standard/age_group_standard_repository.go -destination=mocks/domain/standard/mock_age_group_standard_repository.go -package=mocks
	go tool mockgen -source=internal/domain/standard/rank_standard.go -destination=mocks/domain/standard/mock_rank_standard.go -package=mocks
	go tool mockgen -source=internal/domain/standard/rank_standard_repository.go -destination=mocks/domain/standard/mock_rank_standard_repository.go -package=mocks
	go tool mockgen -source=internal/domain/tenant/profile.go -destination=mocks/domain/tenant/mock_profile.go -package=mocks
	go tool mockgen -source=internal/domain/tenant/profile_repository.go -destination=mocks/domain/tenant/mock_profile_repository.go -package=mocks
	go tool mockgen -source=internal/application/tenant/usecase.go -destination=mocks/application/tenant/mock_usecase.go -package=mocks
	go tool mockgen -source=internal/application/auth/service.go -destination=mocks/application/auth/mock_service.go -package=mocks
	go tool mockgen -source=internal/application/user/service.go -destination=mocks/application/user/mock_service.go -package=mocks
	go tool mockgen -destination=mocks/external/auth/v1/mock_client.go -package=mocks github.com/qkitzero/auth-service/gen/go/auth/v1 AuthServiceClient
	go tool mockgen -destination=mocks/external/group/v1/mock_client.go -package=mocks github.com/qkitzero/user-service/gen/go/group/v1 GroupServiceClient

MIGRATIONS_DIR=internal/infrastructure/db/migrations
MIGRATE=migrate -source file://$(MIGRATIONS_DIR) -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_HOST_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)"

migrate-up:
	$(MIGRATE) up

migrate-up-one:
	$(MIGRATE) up 1

migrate-down:
	$(MIGRATE) down 1

migrate-reset:
	$(MIGRATE) drop -f

migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -format 20060102150405 $(name)

migrate-status:
	$(MIGRATE) version

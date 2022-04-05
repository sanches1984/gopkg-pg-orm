env:
	@cp .env.example ./repository/dao/.env
	@cp .env.example ./migrate/.env
	@docker run --name gopkg-test-db -e POSTGRES_PASSWORD=password -p 4444:5432 -d postgres

test:
	go test -cover -v ./...
build:
	@go build -o ./footballbet-escrow main.go

test:
	@go test -v ./...

run: build
	@./footballbet-escrow

migration:
	@go run migrate/migrate.go
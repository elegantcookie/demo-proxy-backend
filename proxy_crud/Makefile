
start_test_psql:
	docker run -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 postgres

swagger:
	swag init -g ./cmd/app/main.go -o ./docs
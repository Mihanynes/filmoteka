build:
	docker-compose build filmoteka

run:
	docker-compose up filmoteka

migrate:
	migrate -path ./schema -database 'postgres://postgres:111111@0.0.0.0:5436/postgres?sslmode=disable' up
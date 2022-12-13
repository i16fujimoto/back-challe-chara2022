# Makefile
# command: make タグ名

build:
	docker-compose up -d --build
build-dev:
	docker-compose -f docker-compose.dev.yml up -d --build
up-dev:
	docker-compose  -f docker-compose.dev.yml up -d
db-init:
	docker exec mongo mongoimport --host="localhost" --port=27017 --username="root" --password="root" --db="test_database" --collection="test_import" --type="json" --file="./docker-entrypoint-initdb.d/bears.json" --jsonArray
migrate:
	docker-compose exec api go run migrator/migrator.go
dc_db_build:
	docker-compose -f docker-compose.db.yml up --remove-orphans --build

dc_db:
	docker-compose -f docker-compose.db.yml up --remove-orphans

dc_inmem_build:
	docker-compose -f docker-compose.inmem.yml up --remove-orphans --build

dc_inmem:
	docker-compose -f docker-compose.inmem.yml up --remove-orphans

test:
	go test -race -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html
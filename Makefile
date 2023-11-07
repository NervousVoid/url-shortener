dc -storage=db:
	docker run
	docker-compose up  --remove-orphans --build

dc -storage=inmem:
	docker-compose up  --remove-orphans --build

test:
	go test -race -coverprofile=cover.out ./... && go tool cover -html=cover.out -o cover.html
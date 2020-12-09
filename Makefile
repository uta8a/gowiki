dev:
	sudo docker-compose -f ./deployments/docker-compose.dev.yml build
up:
	sudo docker-compose -f ./deployments/docker-compose.dev.yml up
down:
	sudo docker-compose -f ./deployments/docker-compose.dev.yml down
fmt:
	gofmt -w cmd/wiki/main.go
migrate:
	sudo docker-compose -f ./deployments/docker-compose.dev.yml run --rm migration /bin/bash 

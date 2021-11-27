
build:
	docker build -t golang-multiplayer-server:latest ./GameServer

up:
	docker-compose -f ./GameServer/docker-compose.yaml up

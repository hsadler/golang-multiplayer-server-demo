
build:
	docker build -t golang-multiplayer-server:dev -f Dockerfile.dev .

up:
	docker-compose -f docker-compose.dev.yaml up

prod-deploy: heroku-push heroku-release

prod-info:
	heroku apps:info golang-multiplayer-server-demo

prod-logs:
	heroku logs -a golang-multiplayer-server-demo --tail

heroku-push:
	heroku container:push web -a golang-multiplayer-server-demo

heroku-release:
	heroku container:release web -a golang-multiplayer-server-demo


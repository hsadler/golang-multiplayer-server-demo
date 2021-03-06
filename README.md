# Golang Multiplayer Server Demo

A demo game using Unity for the game client and Golang for a websocket server.


### Requirements
- Unity Editor: `version 2020.3.21f1`
- Docker
- Docker Compose
- make


### Running Locally

1. Build and spin-up the Golang game server.
```sh
cd GameServer/
make build
make up
```

2. Open the Unity editor and run the game.


### Deployment

Deploy to Heroku. View app info. View logs.
```sh
cd GameServer/
make prod-deploy
make prod-info
make prod-logs
```

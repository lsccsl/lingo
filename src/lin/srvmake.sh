go mod tidy -compat=1.19

go build -o ../../output/msgque_center lin/server/server_msgque_center
go build -o ../../output/msgque        lin/server/server_msgque
go build -o ../../output/gamesrv       lin/server/server_game
go build -o ../../output/centersrv       lin/server/server_center


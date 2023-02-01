go mod tidy -compat=1.19

go build -o ../../output/msgque_center goserver/server/server_msgque_center
go build -o ../../output/msgque        goserver/server/server_msgque
go build -o ../../output/gamesrv       goserver/server/server_game
go build -o ../../output/centersrv       goserver/server/server_center


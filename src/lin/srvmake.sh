go mod tidy -compat=1.19

go build -o ../../output/msgque_center lin/server_msgque_center
go build -o ../../output/msgque lin/server_msgque


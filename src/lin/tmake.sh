go mod tidy -compat=1.17

go build -o ../../output/testsrv lin/testsrv
go build -o ../../output/testepoll lin/testepoll
go build -o ../../output/test_api lin/test_api

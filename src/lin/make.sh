go mod tidy -compat=1.17
go build -o ../../output/srv lin/srv
go build -o ../../output/epollsrv lin/srv_epoll
go build -o ../../output/testsrv lin/testsrv
go build -o ../../output/testepoll lin/testepoll


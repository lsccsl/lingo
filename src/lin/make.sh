go mod tidy -compat=1.19
go build -o ../../output/srv lin/srv
go build -o ../../output/epollsrv lin/srv_epoll
go build -o ../../output/testsrv lin/testsrv
go build -o ../../output/testepoll lin/testepoll
go build -o ../../output/test_api lin/test_api
go build -o ../../output/test_cross_link lin/test_cross_link
go build -o ../cpp/navwrapper/bin/test_recast lin/test_recast
go build -o ../../output/test_mysql lin/testmysql
go build -o ../../output/test_redis lin/testredis

cp ./msgpacket/msg.pb.go ../../cfg/

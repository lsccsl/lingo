go mod tidy -compat=1.19
go build -o ../../output/srv                  lin/test/srv
go build -o ../../output/epollsrv             lin/test/srv_epoll
go build -o ../../output/testsrv              lin/test/testsrv
go build -o ../../output/testepoll            lin/test/testepoll
go build -o ../../output/test_api             lin/test/test_api
go build -o ../../output/test_cross_link      lin/test/test_cross_link
go build -o ../cpp/navwrapper/bin/test_recast lin/test/test_recast
go build -o ../../output/test_mysql           lin/test/testmysql
go build -o ../../output/test_redis           lin/test/testredis

go build -o ../../output/msgque_center lin/server/server_msgque_center
go build -o ../../output/msgque        lin/server/server_msgque

cp ./msgpacket/*.pb.go ../../cfg/

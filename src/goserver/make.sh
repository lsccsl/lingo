go mod tidy -compat=1.19
go build -o ../../output/srv                  goserver/test/srv
go build -o ../../output/epollsrv             goserver/test/srv_epoll
go build -o ../../output/testsrv              goserver/test/testsrv
go build -o ../../output/testepoll            goserver/test/testepoll
go build -o ../../output/test_api             goserver/test/test_api
go build -o ../../output/test_cross_link      goserver/test/test_cross_link
go build -o ../cpp/navwrapper/bin/test_recast goserver/test/test_recast
go build -o ../../output/test_mysql           goserver/test/testmysql
go build -o ../../output/test_redis           goserver/test/testredis
go build -o ../../output/test_mongo           goserver/test/test_mongo

go build -o ../../output/msgque_center goserver/server/server_msgque_center
go build -o ../../output/msgque        goserver/server/server_msgque
go build -o ../../output/srv_game      goserver/server/server_game
go build -o ../../output/srv_center    goserver/server/server_center
go build -o ../../output/srv_logon     goserver/server/server_logon
go build -o ../../output/srv_db        goserver/server/server_db

cp ./msgpacket/*.pb.go ../../cfg/

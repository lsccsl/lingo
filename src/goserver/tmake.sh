go mod tidy -compat=1.17

go build -o ../../output/testsrv goserver/testsrv
go build -o ../../output/testepoll goserver/testepoll
go build -o ../../output/test_api goserver/test_api
go build -o ../cpp/navwrapper/bin/test_recast goserver/test_recast
go build -o ../../output/test_api goserver/test_cross_link


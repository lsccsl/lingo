protoc --go_out=../msgpacket --proto_path=./ ./*.proto
protoc --cpp_out=../../cpp/test --proto_path=./ ./*.proto
protoc --csharp_out=./ --proto_path=./ ./*.proto

copy /y ..\msgpacket\*.pb.go ..\..\..\cfg\

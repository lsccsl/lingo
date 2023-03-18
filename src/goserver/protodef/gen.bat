protoc --go_out=../msgpacket --proto_path=./ ./*.proto
protoc --cpp_out=../../cpp/test --proto_path=./ ./*.proto
protoc --cpp_out=../../cpp/test_ui --proto_path=./ ./*.proto
protoc --csharp_out=../../../unity_client/test_client/Assets/src --proto_path=./ ./*.proto

copy /y ..\msgpacket\*.pb.go ..\..\..\cfg\

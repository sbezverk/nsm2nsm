package serverendpoints

//go:generate protoc -I . serverendpoints.proto --go_out=plugins=grpc:. --proto_path=$GOPATH/src

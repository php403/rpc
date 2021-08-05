package api

//go:generate protoc -I. -I$GOPATH/src --go_out=plugins=grpc:. --go_opt=paths=source_relative protocol/protocol.proto
//go:generate protoc -I. -I$GOPATH/src --go_out=plugins=grpc:. --go_opt=paths=source_relative comet/comet.proto
//go:generate protoc -I. -I$GOPATH/src --go_out=plugins=grpc:. --go_opt=paths=source_relative logic/logic.proto

//protoc --go_out=./api/logic --proto_path=./api/protocol --proto_path=./api/logic --go_opt=paths=source_relative \
//--go-grpc_out=. --go-grpc_opt=paths=source_relative \
//api/logic/logic.proto
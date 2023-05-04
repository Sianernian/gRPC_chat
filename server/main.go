package main

import (
	pb "gRPC_chat/server/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listenner, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("net.liten err:%v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatRoomServer(grpcServer, &service{})
	if err = grpcServer.Serve(listenner); err != nil {
		log.Fatalf("grpcServer.serve err: %v", err)
	}
}

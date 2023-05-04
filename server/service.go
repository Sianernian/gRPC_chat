package main

import (
	"context"
	"fmt"
	pb "gRPC_chat/server/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

type service struct {
	pb.UnimplementedChatRoomServer
	chatMessageCache []*pb.ChatMessage
	userMap          sync.Map
	L                sync.RWMutex
}

var (
	workers map[pb.ChatRoom_ChatServer]pb.ChatRoom_ChatServer = make(map[pb.ChatRoom_ChatServer]pb.ChatRoom_ChatServer)
)

// 用户注册
func (s *service) Login(ctx context.Context, in *pb.User) (*wrappers.StringValue, error) {
	in.Id = uuid.New().String()
	if _, ok := s.userMap.Load(in.Id); ok {
		return nil, status.Errorf(codes.AlreadyExists, "已有同名用户，请更换")
	}
	s.userMap.Store(in.Id, in)
	go s.Send(nil, &pb.ChatMessage{
		Id:      "server",
		Content: fmt.Sprintf("%v 加入聊天室", in.Name),
		Time:    uint64(time.Now().Unix()),
	})
	return &wrappers.StringValue{Value: in.Id}, status.New(codes.OK, "").Err()
}

func (s *service) Chat(stream pb.ChatRoom_ChatServer) error {
	if s.chatMessageCache == nil {
		s.chatMessageCache = make([]*pb.ChatMessage, 0, 1024)
	}
	workers[stream] = stream
	for _, message := range s.chatMessageCache {
		fmt.Println("message:", message)
		stream.Send(message)
	}
	s.recvMessage(stream)
	return status.New(codes.OK, "").Err()
}

func (s *service) recvMessage(stream pb.ChatRoom_ChatServer) {
	md, _ := metadata.FromIncomingContext(stream.Context())
	for {
		mesg, err := stream.Recv()
		if err != nil {
			s.L.Lock()
			delete(workers, stream)
			s.userMap.Delete(md.Get("uuid")[0])
			fmt.Printf("%s用户掉线,目前用户在线数量:%d", md.Get("uuid")[0], len(workers))
			break
		}
		s.chatMessageCache = append(s.chatMessageCache, mesg)
		v, ok := s.userMap.Load(md.Get("uuid")[0])
		if !ok {
			fmt.Println("用户不存在")
			return
		}

		mesg.Name = v.(*pb.User).Name
		mesg.Time = uint64(time.Now().Unix())
		s.sendMessage(stream, mesg)
		pterm.Info.Println("目前用户在线数量:", len(workers))
	}
}

func (s *service) sendMessage(steam pb.ChatRoom_ChatServer, mes *pb.ChatMessage) {
	s.L.Lock()
	for _, room_chatServer := range workers {
		if room_chatServer != steam {
			err := room_chatServer.Send(mes)
			if err != nil {
				continue
			}
		}
	}
	s.L.Unlock()
}

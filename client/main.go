package main

import (
	"bufio"
	"context"
	pb "gRPC_chat/server/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"strings"
	"time"
)

func main() {
	// 链接服务器
	spinn, _ := pterm.DefaultSpinner.Start("正在链接服务器")
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		spinn.Fail("链接失败")
		pterm.Fatal.Printf("无法连接到服务器:%v", err)
		return
	}
	c := pb.NewChatRoomClient(conn)
	spinn.Success("链接成功")

	// 注册用户名
	var val *wrappers.StringValue
	var user *pb.User
	for {
		rerult, _ := pterm.DefaultInteractiveTextInput.Show("创建用户")
		if strings.TrimSpace(rerult) == "" {
			pterm.Error.Printfln("进入聊天室失败，用户未取名")
			continue
		}
		user = &pb.User{Name: rerult}
		val, err = c.Login(context.TODO(), user)
		if err != nil {
			pterm.Error.Printfln("进入聊天室失败 err:%v", err)
			continue
		} else {
			break
		}
	}
	user.Id = val.Value
	pterm.Success.Println("创建成功！！！")

	// 聊天室逻辑
	stream, _ := c.Chat(metadata.AppendToOutgoingContext(context.Background(), "uuid", user.Id))
	go func(client pb.ChatRoom_ChatClient) {
		for {
			res, err := stream.Recv()
			switch res.Name {
			case "server":
				pterm.Success.Printfln("(%[2]v) [服务器] %[1]s ", res.Content, time.Unix(int64(res.Time), 0).Format(time.ANSIC))
			default:
				pterm.Info.Printfln("(%[3]v) %[1]s : %[2]s", res.Name, res.Content, time.Unix(int64(res.Time), 0).Format(time.ANSIC))
			}
			if err == io.EOF {
				break
			}
		}
	}(stream)
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		input = strings.TrimRight(input, "\r \n")
		stream.Send(&pb.ChatMessage{Id: user.Id, Content: input})
	}
}

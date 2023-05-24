package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	pb "gRPC_chat/server/proto"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/pterm/pterm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

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

		file, err := os.OpenFile("E:\\GoProject\\src\\gRPC_chat/user.json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Fatalf("os.openfile err:%+v\n", err)
		}
		defer file.Close()

		user = &pb.User{Name: rerult}
		user.Id = uuid.New().String()

		Info := User{Id: user.Id, Name: user.Name}

		data, err := json.Marshal(Info)
		if err != nil {
			log.Fatalf("json marshal err:%+v\n", err)
		}

		fmt.Println("data:", string(data))

		IDD, err := file.Write(data)
		if err != nil {
			log.Fatalf("file.Write err:%+v\n", err)
		}
		fmt.Println("file write ID:", IDD)

		val, err = c.Login(context.TODO(), user)
		if err != nil {
			pterm.Error.Printfln("进入聊天室失败 Login err:%v", err)
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

func mergeJSONData(jsonData string, mergedData *[]User) {
	var person User
	err := json.Unmarshal([]byte(jsonData), &person)
	if err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}

	// 将每个person对象添加到mergedData切片中
	*mergedData = append(*mergedData, person)
}

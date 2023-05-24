package main

import (
	"encoding/json"
	"fmt"
	pb "gRPC_chat/server/proto"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
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

func FileModifi() ([]byte, []UserInfo) {
	file, err := os.Open("E:\\GoProject\\src\\gRPC_chat/user.json")
	if err != nil {
		log.Fatalf("os.open err:%+v", err)
	}
	defer file.Close()

	buf := make([]byte, 1024)
	var jsonData []byte
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			break
		}
		jsonData = append(jsonData, buf[:n]...)

	}
	jsonData = []byte(strings.ReplaceAll(string(jsonData), "}]{", "},{"))
	jsonData = []byte("[" + string(jsonData) + "]")
	fmt.Println("jsonData:", string(jsonData))
	count := strings.Count(string(jsonData), "[")
	fmt.Println("Occurrences:", count)
	var result string
	if count > 1 {
		for i := 0; i < count; i++ {
			result = strings.Replace(string(jsonData), "[", "", 1)
			fmt.Println("result11111:", result)
		}
		fmt.Println("result:", result)
		//	result := strings.Replace(string(jsonData), "[", "", -1)
		resultByte := []byte(result)
		//	fmt.Println("resultByte:", string(resultByte))
		err = ioutil.WriteFile("E:\\GoProject\\src\\gRPC_chat/user.json", resultByte, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var info []UserInfo

		err = json.Unmarshal(resultByte, &info)
		if err != nil {
			log.Printf("JSON unmarshal error:%+v", err)
		}

		fmt.Println("return data:", string(resultByte))
		return resultByte, info
	} else {
		err = ioutil.WriteFile("E:\\GoProject\\src\\gRPC_chat/user.json", jsonData, 0644)
		if err != nil {
			log.Fatal(err)
		}
		var info []UserInfo

		err = json.Unmarshal(jsonData, &info)
		if err != nil {
			log.Printf("JSON unmarshal error:%+v", err)
		}

		fmt.Println("return data:", string(jsonData))
		return jsonData, info
	}

}

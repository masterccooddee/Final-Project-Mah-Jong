package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/go-zeromq/zmq4"
)

type Position struct {
	Pos map[string]int
}

var m mao
var RoomID string

func serverexit(conn net.Conn) {
	data := make([]byte, 4096)

	var num int
	for {

		num, _ = conn.Read(data)
		if num == 0 {
			fmt.Println("\nConnection closed")
			os.Exit(1)
		}

		in := string(data[:num])
		in = strings.TrimSpace(in)
		out := strings.Split(in, " ")
		if out[0] == "True" {
			if out[1] == "Room" {
				RoomID = out[2]
			}
		}

		fmt.Printf("\rFrom Server -> %s\n", string(data[:num]))

	}

}
func main() {

	// fmt.Print("Enter server IP: ")
	// var serverip string
	// fmt.Scanln(&serverip)
	// address := net.ParseIP(serverip)
	// if address == nil {
	// 	fmt.Println("Invalid IP")
	// 	goto start
	// }

	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()

	if err != nil {
		fmt.Println("Error dialing", err)
		return
	}
	var dealer zmq4.Socket
	go serverexit(conn)
	var ID string
	for {
		fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		out := strings.Split(text, " ")
		if out[0] == "LOGIN" {
			out[1] = strings.TrimSpace(out[1])
			ID = out[1]
			dealer = zmq4.NewDealer(context.Background(), zmq4.WithID(zmq4.SocketIdentity(out[1])))
			defer dealer.Close()

			err := dealer.Dial("tcp://localhost:7125")
			if err != nil {
				fmt.Println("Error connecting dealer:", err)
				return
			}
			fmt.Println(out[1])

			go func() {

				// DEALER 接收消息
				msg, err := dealer.Recv()
				if err != nil {
					fmt.Println("Error receiving message:", err)
					return
				}
				fmt.Println("Received message:", string(msg.Frames[0]))
				msg, _ = dealer.Recv()
				var pos Position
				json.Unmarshal(msg.Frames[0], &pos)
				fmt.Println(pos.Pos)
				fmt.Println(pos.Pos[ID])

			}()
		}
		if strings.TrimSpace(out[0]) == "CHG" {
			fmt.Println("RoomID: ", RoomID)
			for {
				fmt.Print("Enter text2: ")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				dealer.SendMulti(zmq4.NewMsgFrom([]byte(RoomID), []byte(text)))
			}
		}

		conn.Write([]byte(text))

	}

}

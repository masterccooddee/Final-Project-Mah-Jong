package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-zeromq/zmq4"
)

type player_in struct {
	ID      string
	conn    net.Conn
	Room_ID int
}

var playerlist = make(map[string]*player_in)

var mutex = sync.Mutex{}

func IDcheck(ID string) bool {
	f, err := os.Open("playerlist.txt")
	if err != nil {
		log.Println("Error reading:", err)
		return false
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}
		ID_l := string(line)
		if strings.TrimSpace(ID_l) == strings.TrimSpace(ID) {
			return true
		}

	}
	return false
}

/******************************* Command ******************************\
 * 		ROOM MAKE 			(Create a room)       *
 *		ROOM JOIN XX 			(XX is room ID)       *
 * 		ROOM FIND 			(Auto find a room)    *
 * 		ROOM LEAVE 			(Leave the room)      *
 * 		REG XX 				(Register ID XX:ID)   *
 * 		LOGIN XX 			(Login ID XX:ID)      *
 * 		LOGOUT 				(Logout)              *
\**********************************************************************/

func Cli_handle(conn net.Conn, player player_in, ctx context.Context) {
	// 3. 讀取資料
	defer conn.Close()
	data := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetReadDeadline((time.Now().Add(time.Second * 10)))
			n, err := conn.Read(data)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					// 超時錯誤，繼續檢查上下文
					continue
				}
			}
			if n == 0 {
				if player.ID != "" {
					if player.Room_ID != -1 {
						room := roomlist[player.Room_ID]
						room.leave_room(&player)
					}
					delete(playerlist, player.ID)
				}
				log.Printf("[#FFA500]%v[reset] Connection closed\n", conn.RemoteAddr())
				return
			}
			if err != nil {
				log.Println("[red]ERROR:[reset] ", err)
				return
			}

			data_s := string(data[:n])
			data_s = strings.TrimSpace(data_s)

			command := strings.Split(data_s, " ")
			switch command[0] {
			case "ROOM":
				if player.ID == "" {
					conn.Write([]byte("Please login first\n"))
					continue
				}

				if len(command) < 2 {
					conn.Write([]byte("False Command\n"))
					continue
				}

				switch command[1] {
				case "MAKE":
					if player.Room_ID != -1 {
						conn.Write([]byte("You are already in a room\n"))
						log.Printf("[red]ERROR:[reset] %s try to [yellow]MAKE[reset] room while in a room\n", player.ID)
						continue
					}

					room_id := rand.Intn(99) + 1
					_, exist := roomlist[room_id]
					for exist {
						room_id = rand.Intn(99) + 1
						_, exist = roomlist[room_id]
					}
					makeRoom(room_id, true)
					player.Room_ID = room_id
					log.Printf("[#FFA500]%s[reset] create room %d\n", player.ID, room_id)
					room := roomlist[room_id]
					room.go_in_room(&player, room_id)

				case "JOIN":
					if player.Room_ID != -1 {
						conn.Write([]byte("You are already in a room\n"))
						log.Printf("[red]ERROR:[reset] %s try to [yellow]JOIN[reset] room while in a room\n", player.ID)
						continue
					}
					if len(command) < 3 {
						conn.Write([]byte("False Command\n"))
						continue
					}
					room_id, err := strconv.Atoi(command[2])
					if err != nil {
						conn.Write([]byte("False Command\n"))
						continue
					}
					room := roomlist[room_id]
					room.go_in_room(&player, room_id)

				case "FIND":
					if player.Room_ID != -1 {
						conn.Write([]byte("You are already in a room\n"))
						log.Printf("[red]ERROR:[reset] %s try to [yellow]FIND[reset] room while in a room\n", player.ID)
						continue
					}

					// 等待清空房間
					for cleaning {

					}

					mutex.Lock()
					Room_finder(&player)
					mutex.Unlock()

				case "LEAVE":
					if player.Room_ID == -1 {
						conn.Write([]byte("You are not in a room\n"))
						log.Printf("[red]ERROR:[reset] %s try to [yellow]LEAVE[reset] room while not in a room\n", player.ID)
						continue
					}
					room := roomlist[player.Room_ID]
					room.leave_room(&player)

				default:
					conn.Write([]byte("False Command\n"))
					continue
				}

			case "REG":

				if IDcheck(command[1]) {
					conn.Write([]byte("False ID already exist\n"))
					log.Println("[red]ERROR:[reset] ID already exist")
					continue
				}

				file, err := os.OpenFile("playerlist.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println("[red]ERROR:[reset] ", err)
					continue
				}
				defer file.Close()
				str := command[1] + "\n"
				_, err = file.WriteString(str)
				if err != nil {
					fmt.Println(err)
				}
				player.conn = conn
				player.ID = command[1]
				player.Room_ID = -1
				playerlist[player.ID] = &player
				str_success := "Register success, " + "ID: " + command[1] + "\n"
				conn.Write([]byte(str_success))
				log.Println("Register success, ID: [#FFA500]" + command[1] + "[reset]")

			case "LOGIN":
				if IDcheck(command[1]) == false {
					conn.Write([]byte("False ID not exist\n"))
					log.Println("[red]ERROR:[reset] ID not exist")
					continue
				}

				if _, ok := playerlist[command[1]]; ok {
					conn.Write([]byte("False ID already login\n"))
					log.Println("[red]ERROR:[reset] ID already login")
					continue
				}

				player.conn = conn
				player.ID = command[1]
				player.Room_ID = -1
				playerlist[player.ID] = &player
				conn.Write([]byte("Welcome back " + player.ID + "\n"))
				log.Println("[#FFA500]" + player.ID + "[reset] login")

			case "LOGOUT":
				log.Println("[#FFA500]" + player.ID + "[reset] logout")
				delete(playerlist, player.ID)
				conn.Close()

			}

		}

	}

}

func zmqrecv() {
	for {
		//可能要預防send recv同時進行
		msg, err := router.Recv()
		if err != nil {
			log.Println("Error receiving message:", err)
			return
		}

		player_name := strings.TrimSpace(string(msg.Frames[0]))
		ROOMID := strings.TrimSpace(string(msg.Frames[1]))
		msgout := strings.TrimSpace(string(msg.Frames[2]))

		roomID, err := strconv.Atoi(ROOMID)
		if err != nil {
			log.Println("Error converting ROOMID to int:", err)
			continue
		}
		room := roomlist[roomID]

		msglog := fmt.Sprintf("Room [yellow]%d[reset]: [#FFA500]%s[reset] -> %s", roomID, player_name, msgout)
		zmqloger.Println(msglog)

		room.recvchan <- msgout

	}
}

var cancel context.CancelFunc
var ctx context.Context

func startserver() {

	ctx, cancel = context.WithCancel(context.Background())
	//建立tcp连接
	// 1. 建立監聽器
	ln, err := net.Listen("tcp", ":8080")
	defer ln.Close()
	if err != nil {
		log.Println("Error listening:", err)
		return
	}

	fmt.Fprintln(textView, "Server is listening on port 8080")
	log.SetOutput(textView)
	// ROOM 的規則運行與player的互動
	router = zmq4.NewRouter(context.Background(), zmq4.WithLogger(zmqloger))
	defer router.Close()

	// ROUTER 監聽端點
	err = router.Listen("tcp://*:7125")
	if err != nil {
		log.Fatal("Error starting router:", err)
		return
	}

	fmt.Fprintln(textView, "Router is listening on port 7125")

	go zmqrecv()

	go RoomCleaner(ctx)

	// 2. 建立連線

	for {

		select {
		case <-ctx.Done():
			return
		default:
			conn, err := ln.Accept()
			log.Printf("[#FFA500]%v[reset] Connection established\n", conn.RemoteAddr())
			if err != nil {
				log.Println("Error accepting:", err)
				return

			}
			//conn.Write([]byte("Welcome to the server\n"))
			go Cli_handle(conn, player_in{}, ctx)
		}

	}

}

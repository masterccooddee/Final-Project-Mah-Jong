package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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

func Cli_handle(conn net.Conn, player player_in) {
	// 3. 讀取資料
	defer conn.Close()
	data := make([]byte, 4096)
	for {
		n, err := conn.Read(data)
		if n == 0 {
			fmt.Printf("\n%v Connection closed\n", conn.RemoteAddr())
			return
		}
		if err != nil {
			fmt.Println("Error reading:", err)
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
					continue
				}

				room_id := rand.Intn(99) + 1
				_, exist := roomlist[room_id]
				for exist {
					room_id = rand.Intn(99) + 1
					_, exist = roomlist[room_id]
				}
				makeRoom(room_id, true)
				room := roomlist[room_id]
				room.go_in_room(&player, room_id)

			case "JOIN":
				if player.Room_ID != -1 {
					conn.Write([]byte("You are already in a room\n"))
					continue
				}
				room_id, _ := strconv.Atoi(command[2])
				room := roomlist[room_id]
				room.go_in_room(&player, room_id)

			case "FIND":
				if player.Room_ID != -1 {
					conn.Write([]byte("You are already in a room\n"))
					continue
				}

				mutex.Lock()
				Room_finder(&player)
				mutex.Unlock()

			case "LEAVE":
				if player.Room_ID == -1 {
					conn.Write([]byte("You are not in a room\n"))
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
				continue
			}

			file, err := os.OpenFile("playerlist.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Println("Error reading:", err)
				continue
			}
			defer file.Close()
			str := command[1] + "\n"
			_, err = file.WriteString(str)
			if err != nil {
				fmt.Println(err)
			}
			playerlist[command[1]] = &player_in{ID: command[1], conn: conn, Room_ID: -1}
			player.ID = command[1]
			player.conn = conn
			player.Room_ID = -1
			str_success := "Register success, " + "ID: " + command[1] + "\n"
			conn.Write([]byte(str_success))

		case "LOGIN":
			if IDcheck(command[1]) == false {
				conn.Write([]byte("False ID not exist\n"))
				continue
			}

			player.conn = conn
			player.ID = command[1]
			player.Room_ID = -1
			playerlist[player.ID] = &player
			conn.Write([]byte("Welcome back " + player.ID + "\n"))

		case "LOGOUT":
			delete(playerlist, player.ID)
			conn.Close()

		}

		log.Println(playerlist)

	}

}

func startserver() {
	//建立tcp连接
	// 1. 建立監聽器
	ln, err := net.Listen("tcp", ":8080")
	defer ln.Close()
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	fmt.Println("Server is listening on port 8080")
	// 2. 建立連線

	for {
		conn, err := ln.Accept()
		fmt.Printf("%v Connection established\n", conn.RemoteAddr())
		if err != nil {
			fmt.Println("Error accepting:", err)
			return

		}
		conn.Write([]byte("Welcome to the server\n"))
		go Cli_handle(conn, player_in{})

	}

}
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"
)

var conn net.Conn
var err error
var RoomID string

func rrecv() string {

	data := make([]byte, 4096)

	var num int

	num, _ = conn.Read(data)
	if num == 0 {
		fmt.Println("\nConnection closed")
		os.Exit(1)
	}

	in := string(data[:num])
	in = strings.TrimSpace(in)
	return in
}

func TestRest(t *testing.T) {
	conn, err = net.Dial("tcp", "localhost:8080")
	defer conn.Close()

	if err != nil {
		log.Fatal("connect error")
	}

	connectTest := []struct {
		expect      string
		testcommand string
	}{
		{"Please login first", "ROOM"},
		{"False Command", "REG"},
		{"False Command", "REGG"},
		{"False ID already exist", "REG hehehe"},
		{"Register success, ID: 222", "REG 222"},
		{"False Command", "ROOM"},
		{"True Room", "ROOM MAKE"}, //後面還要接數字
		{"You are already in a room", "ROOM FIND"},
		{"You are already in a room", "ROOM JOIN"},
		{"True Leave room", "ROOM LEAVE"},
		{"You are not in a room", "ROOM LEAVE"},
		{"False Command", "ROOM JOIN"},
		{"False Command", "ROOM JOIN ***"},
		{"True Room", "ROOM FIND"}, //後面還要接數字
		{"False ID not exist", "LOGIN 111"},
		{"False Command", "LOGI"},
		{"False Command", "LOGINN"},
		{"Welcome back hehehe", "LOGIN hehehe"},
		{"False ID already login", "LOGIN hehehe"},
		{"False Command", "ROOM"},
		{"True Room", "ROOM MAKE"}, //後面還要接數字
		{"You are already in a room", "ROOM FIND"},
		{"You are already in a room", "ROOM JOIN"},
		{"True Leave room", "ROOM LEAVE"},
		{"You are not in a room", "ROOM LEAVE"},
		{"False Command", "ROOM JOIN"},
		{"False Command", "ROOM JOIN ***"},
		{"True Room", "ROOM FIND"}, //後面還要接數字
	}

	for _, e := range connectTest {
		conn.Write([]byte(e.testcommand))
		t.Log(e.testcommand)
		get := rrecv()
		if get != e.expect {
			t.Errorf("Failed,\nExpected to be %s\nActual : %s", e.expect, get)
		}
	}

}

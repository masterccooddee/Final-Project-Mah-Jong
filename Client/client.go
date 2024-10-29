package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type mao struct {
	Card []string
	Wan  []string
	Tong []string
	Tiao []string
}

var m mao

func serverexit(conn net.Conn) {
	data := make([]byte, 4096)

	var num int
	for {

		num, _ = conn.Read(data)
		if num == 0 {
			fmt.Println("\nConnection closed")
			os.Exit(1)
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
	go serverexit(conn)
	for {
		fmt.Print("Enter text: ")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		conn.Write([]byte(text))

	}

}

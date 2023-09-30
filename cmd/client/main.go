package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("ping"))
	reader := bufio.NewReader(conn)
	line, _, _ := reader.ReadLine()
	fmt.Println(string(line))
}

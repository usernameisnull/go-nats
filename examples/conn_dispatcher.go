package main

import (
	"net"
	"log"
	"fmt"
	"bufio"
	//"os"
)

func main(){
	tcpAdd := "127.0.0.1:9999"
	conn, err := net.Dial("tcp", tcpAdd)
	if err != nil {
		log.Fatalf("connect to %s, errors: %s\n", tcpAdd, err.Error())
	}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	conn.Write([]byte("$SUB_ALL"))
	for{
		//// read in input from stdin
		//reader := bufio.NewReader(os.Stdin)
		//fmt.Print("Text to send: ")
		//text, _ := reader.ReadString('\n')
		//// send to socket
		//fmt.Fprintf(conn, text + "\n")
		//// listen for reply
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}
}

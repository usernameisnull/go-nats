package main

import (
	"net"
	"log"
	"fmt"
	//"bufio"
	//"os"
)

func main(){
	tcpAdd := "127.0.0.1:9999"
	conn, err := net.Dial("tcp", tcpAdd)
	if err != nil {
		log.Fatalf("connect to %s, errors: %s\n", tcpAdd, err.Error())
	}
	conn.Write([]byte("$SUB_ALL"))
	done := make(chan string)
	//for{
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: "+message)
	//}
	for {
		go readFrDispatcher(conn,done)
		fmt.Println(<-done)
	}

}


func readFrDispatcher(conn net.Conn, done chan string){
	buf := make([]byte, 1024)
	reqLen, _ := conn.Read(buf)
	done<-string(buf[:reqLen-1])
}

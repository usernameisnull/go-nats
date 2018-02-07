// Copyright 2012-2016 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"flag"
	"log"
	"time"
	//"fmt"
	"strings"
	"os"
	"net"
	//"bufio"
	"github.com/nats-io/go-nats"
)

// NOTE: Use tls scheme for TLS, e.g. nats-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

var urls = "nats://localhost:4222"
func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	//log.SetFlags(0)
	//flag.Usage = usage
	//flag.Parse()

	//natsConn, err := nats.Connect(*urls)
	natsConn, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer natsConn.Close()

	subj := os.Args[1]
	//set := os.Args[2]
	beginTime := time.Now()
	printOutStamp := beginTime
	//beginTimeFormatted:=beginTime.Format("2006-01-02 15:04:05")

	var count float64 = 1
	tcpAdd := "127.0.0.1:9999"
	tcpConn, err := net.Dial("tcp", tcpAdd)
	if err != nil {
		log.Fatalf("connect to %s, errors: %s\n", tcpAdd, err.Error())
	}
	tcpConn.Write([]byte("$SUB_ALL"))
	readFrDisDone := make(chan []byte)
	//readFrNats := make(chan []byte)
	//for{
	//
	//	go readFrDispatcher(tcpConn, readFrDisDone)
	//	msg := <-readFrDisDone
	//	err = natsConn.Publish(subj, msg)
	//	natsConn.Flush()
	//	if err != nil {
	//		log.Fatalf("Error during publish: %v\n", err)
	//	}
	//	if time.Since(printOutStamp).Seconds() >= 10{
	//		log.Printf("Published [%s] : '%s'\n", subj, msg)
	//		printOutStamp=time.Now()
	//	}
	//	count++
	//}
	go readFrDispatcher(tcpConn, readFrDisDone)
	go subFromNATS(natsConn, tcpConn)
	for{
		msg:=<-readFrDisDone
		count = count+float64(strings.Count(string(msg),"$POS"))
		//fmt.Println(string(msg))
		err = natsConn.Publish(subj, msg)
		natsConn.Flush()
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		if time.Since(printOutStamp).Seconds() >= 10{
			//log.Printf("Published [%s] : '%s'\n", subj, msg)
			log.Printf("Speed is: %v", count/time.Since(beginTime).Seconds())
			printOutStamp=time.Now()
		}
		count++
		 //select{
			//case msg:=<-readFrDisDone:
			//	fmt.Println(string(msg))
			//	err = natsConn.Publish(subj, msg)
			//	natsConn.Flush()
			//	if err != nil {
			//		log.Fatalf("Error during publish: %v\n", err)
			//	}
			//	if time.Since(printOutStamp).Seconds() >= 10{
			//		log.Printf("Published [%s] : '%s'\n", subj, msg)
			//		printOutStamp=time.Now()
			//	}
			//	count++
			//case m:=<-readFrNats:
			//	tcpConn.Write(m)
			//	natsConn.Flush()
		 //}
	}

	if err := natsConn.LastError(); err != nil {
		log.Fatal(err)
	}
}

// 从dispatcher的tcp端口读取消息，这个消息被pub出去
func readFrDispatcher(conn net.Conn, rfDisDone chan []byte){
	for{
		buf := make([]byte, 1024)
		reqLen, err := conn.Read(buf)
		if reqLen > 1 && err==nil{
			rfDisDone<-buf[:reqLen-1]
			buf = nil
		}
	}
}

// 从消息服务器订阅消息，这个消息被作为指令发送给dispather，用以控制dispatcher的行为。
func subFromNATS(natsCn *nats.Conn, conn net.Conn){
		natsCn.Subscribe("goo", func(msg *nats.Msg) {
			// todo: 向dispatcher发起连接，并发送指令，需要取得与dispatcher连接的conn，tcpConn.Write([]byte("$SUB_ALL"))
			log.Printf("######################## %s\n", string(msg.Data))
			conn.Write(msg.Data)
		})
}


func getSubj() string{
	return os.Args[1]
}

// Copyright 2012-2016 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"flag"
	"log"
	"time"
	//"fmt"
	"os"
	"net"
	//"bufio"
	"github.com/nats-io/go-nats"
	"strings"
)

// NOTE: Use tls scheme for TLS, e.g. nats-pub -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-pub [-s server (%s)] <subject> <msg> \n", nats.DefaultURL)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	subj := os.Args[1]
	//set := os.Args[2]
	beginTime := time.Now()
	printOutStamp := beginTime
	//beginTimeFormatted:=beginTime.Format("2006-01-02 15:04:05")

	var count float64 = 1
	tcpAdd := "127.0.0.1:9999"
	conn, err := net.Dial("tcp", tcpAdd)
	if err != nil {
		log.Fatalf("connect to %s, errors: %s\n", tcpAdd, err.Error())
	}
	conn.Write([]byte("$SUB_ALL"))
	done := make(chan []byte)

	for{
		//msg := fmt.Sprintf("%s %s %v %s %v %s speed = %v",
		//	beginTimeFormatted,
		//		set,
		//	time.Now().Format("2006-01-02 15:04:05"),
		//		set,
		//	count,
		//		set,
		//	count/(time.Since(beginTime).Seconds()))
		go readFrDispatcher(conn, done)
		msg := <-done
		count = count+float64(strings.Count(string(msg),"$POS"))
		err = nc.Publish(subj, msg)
		nc.Flush()
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		if time.Since(printOutStamp).Seconds() >= 10{
			log.Printf("Speed is: %v", count/time.Since(beginTime).Seconds())
			//log.Printf("Published [%s] : '%s'\n", subj, msg)
			printOutStamp=time.Now()
		}
		//count++
	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}

func readFrDispatcher(conn net.Conn, done chan []byte){
	buf := make([]byte, 1024)
	reqLen, _ := conn.Read(buf)
	done<-buf[:reqLen-1]
}

// 从消息服务器订阅消息
func subFromNATS(){

}
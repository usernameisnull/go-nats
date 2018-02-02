// Copyright 2012-2016 Apcera Inc. All rights reserved.
// +build ignore

package main

import (
	"flag"
	"log"
	"time"
	"fmt"
	"os"
	"github.com/nats-io/go-nats"
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
	set := os.Args[2]
	beginTime := time.Now()
	printOutStamp := beginTime
	beginTimeFormatted:=beginTime.Format("2006-01-02 15:04:05")
	var count float64 = 1
	for{
		msg := fmt.Sprintf("%s %s %v %s %v %s speed = %v",
			beginTimeFormatted,
				set,
			time.Now().Format("2006-01-02 15:04:05"),
				set,
			count,
				set,
			count/(time.Since(beginTime).Seconds()))
		err = nc.Publish(subj, []byte(msg))
		nc.Flush()
		if err != nil {
			log.Fatalf("Error during publish: %v\n", err)
		}
		if time.Since(printOutStamp).Seconds() >= 10{
			log.Printf("Published [%s] : '%s'\n", subj, msg)
			printOutStamp=time.Now()
		}
		count++
	}

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
}

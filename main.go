/*
 * @Date: 2020-08-14 14:39:08
 * @LastEditors: liuruihao@huobi.com
 * @LastEditTime: 2020-08-14 17:01:31
 * @FilePath: /tcpproxy/main.go
 */
package main

import (
	"flag"
	"net"

	logging "github.com/ipfs/go-log"
)

var (
	log = logging.Logger("main")
)

func main() {
	linkAddr := flag.String("linkAddr", "", "127.0.0.1:3333")
	bindAddr := flag.String("bindAddr", "", "0.0.0.0:4444")
	levelStr := flag.String("level", "info", "debug,info,warn,error")
	logPath := flag.String("logpath", "./", "./log/")
	flag.Parse()

	logconf := logging.Config{Format: logging.ColorizedOutput, Stderr: false, Stdout: true}
	if *logPath != "" {
		logconf.File = *logPath + "proxy.log"
	}
	if level, err := logging.LevelFromString(*levelStr); err != nil {
		log.Fatal(err)
	} else {
		logconf.Level = level
	}
	logging.SetupLogging(logconf)

	// if *linkAddr == "" {
	// 	log.Fatal("linkaddr is null")
	// }

	link, err := net.ResolveTCPAddr("tcp", *linkAddr)
	if err != nil {
		log.Fatal("linkAddr Resolve failed: ", err)
	}

	bind, err := net.ResolveTCPAddr("tcp", *bindAddr)
	if err != nil {
		log.Fatal("bindAddr Resolve failed: ", err)
	}
	listener, err := net.ListenTCP("tcp", bind)
	if err != nil {
		log.Fatal("ListenTCP failed: ", err)
	}
	log.Info("listenning ", bind.String(), " to ", link.String())
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("accept failed: ", err)
			continue
		}
		go NewChannel(link, tcpConn).Run()
	}
}

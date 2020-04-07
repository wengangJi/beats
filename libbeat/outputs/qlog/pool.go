package qlog

import (
	"fmt"
	"log"
	"net"
	"time"
)

var effectiveConns chan net.Conn
var retryConns chan string

var globalTimeOut time.Duration

func InitConnectPool(servList []string, timeOut time.Duration, minConnNum int) {
	globalTimeOut = timeOut
	effectiveConns = make(chan net.Conn, len(servList)*minConnNum)
	retryConns = make(chan string, len(servList)*minConnNum)
	go Retry()
	for _, addr := range servList {
		connCount := 0
		for connCount < minConnNum {
			connCount = connCount + 1
			connectServer(addr, globalTimeOut)
		}
	}
}

func connectServer(addr string, globalTimeOut time.Duration) {
	conn, err := net.DialTimeout("tcp", addr, globalTimeOut*time.Second)
	if err != nil {
		log.Println("error connecting:", err)
		retryConns <- addr
	} else {
		log.Println("connect to ", conn.RemoteAddr())
		effectiveConns <- conn
	}
}

func Get() net.Conn {
CREATECONN:
	conn := <-effectiveConns
	if conn == nil {
		goto CREATECONN
	}
	return conn
}

func Put(conn net.Conn) {
	effectiveConns <- conn
}

func Drop(conn net.Conn) {
	ipaddr := conn.RemoteAddr()
	conn.Close()
	retryConns <- fmt.Sprintf("%v", ipaddr)
}

func Retry() {
	for addr := range retryConns {
		connectServer(addr, globalTimeOut)
		time.Sleep(5 * time.Second)
	}
}

func Close() {
	for conn := range effectiveConns {
		conn.Close()
	}
	close(effectiveConns)
	close(retryConns)
}

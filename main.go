package main

import (
	"io"
	"log"
	"net"
	"sync"
)

const LocalAddr = ":28001"
const RemoteAddr = "10.20.121.186:28001"

func main() {
	log.Printf("[ TCP Proxy ] local%s -> %s\n", LocalAddr, RemoteAddr)

	lis, e := net.Listen("tcp", LocalAddr)
	if e != nil {
		println(e.Error())
	}
	defer func() {
		_ = lis.Close()
	}()

	// 转发机
	for {
		conn, e := lis.Accept()
		if e != nil {
			log.Printf("[ ERROR ] 建立连接错误: %s\n", e.Error())
			continue
		}
		go handle(conn, RemoteAddr)
	}
}

func handle(conn net.Conn, dialIp string) {
	var dconn net.Conn
	var e error

	defer func() {
		_ = dconn.Close()
		_ = conn.Close()
	}()

	// 和目标地址建立连接
	dconn, e = net.Dial("tcp", dialIp)
	if e != nil {
		log.Printf("[ ERROR ] 和目标地址建立连接: %s\n", e.Error())
		return
	}

	// Log
	log.Printf("[ TCP ] %s -> %s\n", dconn.LocalAddr(), dconn.RemoteAddr())

	// 数据交换
	dataExchange(conn, dconn, dialIp)
}

func dataExchange(conn net.Conn, dconn net.Conn, dialIp string) {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(2)

	go func(conn net.Conn, dconn net.Conn) {
		defer wg.Done()
		if _, e := io.Copy(dconn, conn); e != nil {
			log.Printf("往 %s 发送数据失败:%s\n", dialIp, e.Error())
			return
		}
	}(conn, dconn)

	go func(conn net.Conn, dconn net.Conn) {
		defer wg.Done()
		if _, e := io.Copy(conn, dconn); e != nil {
			log.Printf("从 %s 接收数据失败:%s\n", dialIp, e.Error())
			return
		}
	}(conn, dconn)

}

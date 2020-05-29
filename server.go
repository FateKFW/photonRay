package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

//服务端
type Server struct {
	connType string
	ip string
	port string
	listen net.Listener
}

func bindServerParam() *Server {
	var help bool
	var port string
	flag.BoolVar(&help, "h", false, "Parameter Description")
	flag.StringVar(&port,"p", "8282", "Server port")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	server := &Server{"tcp", "127.0.0.1", port,nil}
	showBanner()
	return server
}

//初始化转发服务
func (server *Server) initServer() error{
	listener, err := net.Listen(server.connType, server.ip + ":" + server.port)
	if err != nil {
		return err
	}
	server.listen = listener
	return nil
}

//开始监听服务
func (server *Server) startServer() {
	log.Print("Service started successfully")
	//延迟关闭
	defer server.listen.Close()
	//接收数据
	for {
		conn, err := server.listen.Accept()
		if err!=nil {
			log.Print(err)
			break
		}
		go handleRequest(conn, server.ip)
	}
}

//处理请求
func handleRequest(conn net.Conn, ip string) {
	defer conn.Close()
	//读取初始连接请求
	buf := bufio.NewReader(conn)
	line, _ , err := buf.ReadLine()
	if err != nil {
		log.Println(err)
	}
	if "PING" == string(line) {
		port := pc.getPort()
		log.Printf("The client connection is successful, open the access port %v\n", port)
		conn.Write([]byte("PONG\nPlease access "+strconv.Itoa(port)))
		//启动接收服务
		listen, err := net.Listen("tcp", ip+":"+strconv.Itoa(port))
		if err != nil {
			log.Println(err)
			return
		}
		for {
			recConn,err := listen.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Request: RemoteAddr %v LocalAddr %v", recConn.RemoteAddr(), recConn.LocalAddr())
			go io.Copy(conn, recConn)
			io.Copy(recConn, conn)
			recConn.Close()
		}
	}
}

func main() {
	initPortConfig()
	var server = bindServerParam()
	err := server.initServer()
	checkErr(err)
	server.startServer()
}

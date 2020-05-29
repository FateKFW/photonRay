package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

type Client struct {
	connType string
	remoteIP string
	remotePort string
	localIP string
	localPort string
	conn net.Conn
}

//参数绑定
func bindClientParam() *Client {
	var help bool
	var sip, sp, ip, port string
	flag.BoolVar(&help, "h", false, "Parameter Description")
	flag.StringVar(&sip,"rs", "127.0.0.1", "Server ip")
	flag.StringVar(&sp,"rp", "8282", "Server port")
	flag.StringVar(&port,"cp", "", "Proxy port")
	flag.StringVar(&ip,"cs", "127.0.0.1", "Proxy address")
	flag.Parse()
	if help || port == "" {
		flag.PrintDefaults()
		time.Sleep(10 * time.Second)
		os.Exit(1)
	}
	cli := &Client{"tcp", sip, sp, ip, port, nil}
	showBanner()
	return cli
}

//初始化客户端连接
func (cli *Client) initClient() error{
	//访问服务器
	client,err := net.ResolveTCPAddr(cli.connType, cli.remoteIP + ":" + cli.remotePort)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP(cli.connType,nil, client)
	if err != nil {
		return err
	}
	cli.conn = conn
	return nil
}

//开启客户端，与服务器通讯
func (cli *Client) startClient() {

	//延迟关闭连接
	defer cli.conn.Close()

	//请求连接，如服务器不通或者数据报文非法则关闭客户端
	cli.conn.Write([]byte("PING\n"))

	res, err := readString(cli.conn,-1)
	checkErr(err)

	if "PONG" == res[:4] {	//通过服务器验证
		log.Println("Server connection successful")
		//回显访问链接
		log.Printf(res[5:])
		for {
			
			localClient,err := net.ResolveTCPAddr(cli.connType, cli.localIP+":"+cli.localPort)
			checkErr(err)
			localConn, err := net.DialTCP(cli.connType,nil, localClient)
			checkErr(err)

			log.Printf("Request: RemoteAddr %v LocalAddr %v", cli.conn.RemoteAddr(), cli.conn.LocalAddr())
			go io.Copy(localConn, cli.conn)
			io.Copy(cli.conn, localConn)

			log.Println("handleOK")

		}
	} else {	//未通过服务器验证
		log.Println("Server connection failed")
		os.Exit(1)
	}
}

func main() {
	var cli = bindClientParam()
	err := cli.initClient()
	checkErr(err)
	cli.startClient()
}
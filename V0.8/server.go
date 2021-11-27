package main

import (
	"fmt"
	"net"
	"sync"
	"io"
	"time"
)

type Server struct {
	Ip   string
	Port int

	//OnlineMap 在线的用户
	OnlineMap map[string]*User

	//广播的列表
	Message chan string

	//一些操作需要加锁
	mapLock sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

//广播消息
func (s *Server) BroadCast(user *User, msg string) {
	sendMes := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- sendMes
}

//监听Message channel 的goroutine，一旦有消息就发送给全部的User
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		//将msg发送给全体User
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.Ch <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//业务
	// fmt.Println("创建连接成功！")
	user := NewUser(conn,s)
	//用户上线,将用户添加进OnlineMap
	user.Online()

	isAction := make(chan bool)

	go func(){
		buf := make([]byte,4096)
		//这边一定要一个死循环监听用户，不然用户只能发送一次消息
		for{
			n,err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Sprintf("read is err:", err)
				return 
			}
			//去除最后一个回车符号
			msg := string(buf[:n-1])
			//广播
			user.DoMessage(msg)

			//代表当前用户是活跃的
			isAction <-true
		}
	}()

	//当前handle阻塞
	for{
		select {
			//判断是否处于活跃，处于活跃则要重置监听时间
		case <-isAction:
			//监听时间
		case <-time.After(time.Second * 1000):
			//发生被踢消息
			user.SendMsg("您被踢了!")

			//释放资源
			close(user.Ch)

			//关闭连接
			conn.Close()
			//关闭当前处理 或者 runtime.Goexit()
			return 
		}
	}
	
}

func (s *Server) Start() {
	//listen
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("listen to err:", err)
		return
	}

	

	//close
	defer listen.Close()

	//启动监听广播msg的goroutine
	go s.ListenMessage()

	for {
		//accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept to err:", err)
			continue
		}

		//to handler
		go s.Handler(conn)
	}

}

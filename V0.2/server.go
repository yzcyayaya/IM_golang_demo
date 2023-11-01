package main

import(
	"fmt"
	"net"
	"sync"	
)

type Server struct{
	Ip string
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
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
	return server
}
//广播消息
func (s *Server)BroadCast(user *User,msg string) {
	sendMes := "[" + user.Name + "]" + user.Addr + ":" + msg;
	s.Message <- sendMes
}

//监听Message channel 的goroutine，一旦有消息就发送给全部的User
func (s *Server)ListenMessage(){
	for{
		
		msg := <-s.Message

		//将msg发送给全体User
		s.mapLock.Lock()
		for _,user := range s.OnlineMap{
			user.Ch <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server)Handler(conn net.Conn) {
	//业务
	// fmt.Println("创建连接成功！")
	user := NewUser(conn)
	//用户上线,将用户添加进OnlineMap
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播当前用户上线信息
	s.BroadCast(user,"OnLine")

	//阻塞
	select{}
}

func (s *Server)Start(){
	//listen
	listen,err := net.Listen("tcp",fmt.Sprintf("%s:%d",s.Ip,s.Port))
	if err != nil{
		fmt.Println("listen to err:",err)
		return
	}

	//close
	defer listen.Close()

	go s.ListenMessage()

	for{
		//accept
		conn,err := listen.Accept()
		if err != nil{
			fmt.Println("accept to err:",err)
			continue
		}

		//to handler
		go s.Handler(conn)
	}
}
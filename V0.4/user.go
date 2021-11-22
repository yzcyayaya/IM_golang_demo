package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	Ch   chan string
	conn net.Conn
	server *Server
}

//上线模块

func (u *User) Online() {
	//用户上线,将用户添加进OnlineMap
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播当前用户上线信息
	u.server.BroadCast(u, "已上线")
}

//下线模块
func (u *User) Offline() {
	//用户下线,将用户添加进OnlineMap
	u.server.mapLock.Lock() 
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()

	//广播当前用户上线信息
	u.server.BroadCast(u, "已下线")
}

//发送消息业务模块
func (u *User) DoMessage(msg string) {
	u.server.BroadCast(u,msg)
}

//监听channel的,一旦有消息直接发送给客户端
func (u *User) LintenMessage() {
	for {
		msg := <-u.Ch
		u.conn.Write([]byte(msg+"\n"))
	}
}

//创建一个用户的API
func NewUser(conn net.Conn,server *Server) *User {
	userAdder := conn.RemoteAddr().String()
	user := &User{
		Name: userAdder,
		Addr: userAdder,
		Ch:   make(chan string),
		conn: conn,
		server: server,
	}
	//监听
	go user.LintenMessage()

	return user
}

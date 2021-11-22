package main

import (
	"net"
	"strings"
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
	u.server.BroadCast(u, "下线")
}
func (u *User) SendMsg(msg string){
	u.conn.Write([]byte(msg))
}

//发送消息业务模块
func (u *User) DoMessage(msg string) {
	if msg == "who"{
		//查询当前在线有哪些
		u.server.mapLock.Lock()
		for _,user := range u.server.OnlineMap {
			onlineMsg := "["+user.Addr+"]"+ user.Name + ":在线...\n"
			u.SendMsg(onlineMsg)
		}
		u.server.mapLock.Unlock()
	}else if len(msg) > 7 && msg[:7] == "rename:"{
		//消息格式 rename:张三
		newName := strings.Split(msg,":")[1]
		
		_,ok := u.server.OnlineMap[newName]
		//判断name是否存在
		if ok{
			u.SendMsg("该用户名已经存在!")
		}else{
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap,u.Name)
			u.server.OnlineMap[newName] = u
			u.server.mapLock.Unlock()
			
			u.Name = newName
			u.SendMsg("您已更新用户名:"+ u.Name + "\n")
		}

	} else {
		u.server.BroadCast(u,msg)
	}
}

//监听channel的,一旦有消息直接发送给客户端
func (u *User) LintenMessage() {
	for {
		msg := <-u.Ch
		u.SendMsg(msg+"\n")
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

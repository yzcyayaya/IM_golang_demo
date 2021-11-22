package main

import (
	"net"
)

type User struct{
	Name string
	Addr string
	Ch chan string
	conn net.Conn	
}
//创建一个用户的API
func NewUser(conn net.Conn) *User{
	userAdder := conn.RemoteAddr().String()
	user := &User{
		Name: userAdder,
		Addr: userAdder,
		Ch: make(chan string),
		conn: conn,
	}
	//监听
	go user.LintenMessage()

	return user
}
//监听channel的,一旦有消息直接发送给客户端
func (u *User)LintenMessage() {
	for{
		msg := <-u.Ch
		u.conn.Write([]byte(msg + "\n"))
	}
}

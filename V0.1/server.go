package main

import(
	"fmt"
	"net"	
)

type Server struct{
	Ip string
	Port int
}


func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip: ip,
		Port: port,
	}
	return server
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
func (s *Server)Handler(conn net.Conn) {
	//业务
	fmt.Println("创建连接成功！")
}
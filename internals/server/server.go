package server

import "net"
import "fmt"
import "dumbbelldog/internals/client"

type Server struct{
  port int
  address string
  paths map[string]client.Callback
}


func NewServer(address string,port int) *Server{
  return &Server{
    port:port,
    address:address,
    paths:make(map[string]client.Callback),
  }
}

func (server *Server) AddPath(path string,callback client.Callback){
  server.paths[path] = callback
}

func (server *Server) Run(){
  tcpAddr, err := net.ResolveTCPAddr("tcp4",fmt.Sprintf("%s:%d",server.address,server.port))
  if err != nil {
    fmt.Println(err)
  }
  listener,err := net.ListenTCP("tcp4",tcpAddr)
  if err != nil {
    fmt.Println(err)
  }
  for {
    clientConn,err := listener.Accept()
    if err != nil {
      fmt.Println(err)
      continue
    }
    client.NewClient(clientConn,server.paths)
  }
}

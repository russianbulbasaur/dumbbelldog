package dumbbelldog

import (
	"dumbbelldog/internals/client"
	"log"
	"net"
)

func main() {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4221")
	if err != nil {
		log.Fatalln(err)
	}
	server, err := net.ListenTCP("tcp", address)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		clientConn, err := server.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		client.NewClient(clientConn)
	}
}

package client

import (
	"log"
	"net"
)
import "context"

type Client struct {
	conn          net.Conn
	clientContext context.Context
	cancelFunc    context.CancelFunc
	writerPipe    chan []byte
}

func NewClient(conn net.Conn) *Client {
	clientContext, cancel := context.WithCancel(context.Background())
	client := Client{
		conn,
		clientContext,
		cancel,
		make(chan []byte),
	}
	go client.forkReader()
	go client.forkWriter()
	return &client
}

func (client *Client) forkReader() {
	var buffer []byte = make([]byte, 1024)
	for {
		n, err := client.conn.Read(buffer)
		if err != nil {
			log.Println(err)
			client.cancelFunc()
			return
		}
        response := client.parseRequest(buffer[0:n])
        client.writerPipe <- response
	}
}

func (client *Client) forkWriter() {
	for {
		select {
		case <-client.clientContext.Done():
			return
		case message := <-client.writerPipe:
			_, err := client.conn.Write(message)
			if err != nil {
				log.Println("Error writing : ", err)
			}
		}
	}
}

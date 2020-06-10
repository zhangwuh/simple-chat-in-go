package main

import (
	"time"

	"github.com/zhangwuh/simple-chat-in-go/cs"
)



func main() {
	shutdown := make(chan bool)
	server := cs.NewDummyServer()
	conCh := make(chan cs.Connection)
	startServer(server, conCh)
	registerClients(server)
	<-shutdown
}

func registerClients(server *cs.DummyServer) {
	client1 := cs.NewDummyClient(cs.NewDummyConnection("client1", make(chan cs.Message)))
	client2 := cs.NewDummyClient(cs.NewDummyConnection("client2", make(chan cs.Message)))
	server.Accept(client1.Connection())
	server.Accept(client2.Connection())
}


func startServer(server cs.Server, ch <-chan cs.Connection) {
	pingClients(server)
}

func pingClients(server cs.Server) {
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				for _, c := range server.Connections() {
					go c.Write([]byte("hello from server"))
				}
			}
		}
	}()

}

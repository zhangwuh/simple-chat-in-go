package cs

import (
	"fmt"
)

type Connection interface {
	Id() string
	Read() chan Message
	Write(msg Message)
	Close()
}

type Message struct {
	content []byte
	target  string
}

func NewMessage(content []byte, target string) Message {
	return Message{content, target}
}

type DummyConnection struct {
	id  string
	ch  chan Message
	seq int
}

func NewDummyConnection(id string, ch chan Message) *DummyConnection {
	return &DummyConnection{id: id, ch: ch}
}

func (con *DummyConnection) Id() string {
	return con.id
}

func (con *DummyConnection) Read() chan Message {
	return con.ch
}

func (con *DummyConnection) Write(message Message) {
	msg := NewMessage(message.content, "")
	defer func() {
		con.seq++
	}()
	con.ch <- msg
}

func (con *DummyConnection) Close() {
	close(con.ch)
	fmt.Println(fmt.Sprintf("connection %s closed", con.Id()))
}

type Server interface {
	Connections() []Connection
	GetConnection(id string) Connection
	Accept(con Connection)
	Close(con Connection)
	Receive(msg Message)
}

type DummyServer struct {
	connections map[string]Connection
	ch          chan Message
}

func NewDummyServer() *DummyServer {
	return &DummyServer{connections: map[string]Connection{}}
}

func (ds *DummyServer) Receive(msg Message) {
	for _, c := range ds.connections {
		if c.Id() == msg.target {
			c.Write(msg)
		}
	}
}

func (ds *DummyServer) Connections() []Connection {
	var conns []Connection
	for _, v := range ds.connections {
		conns = append(conns, v)
	}
	return conns
}

func (ds *DummyServer) Accept(con Connection) {
	ds.connections[con.Id()] = con
}

func (ds *DummyServer) GetConnection(id string) Connection {
	return ds.connections[id]
}

func (ds *DummyServer) Close(con Connection) {
	defer func() {
		if ds.GetConnection(con.Id()) != nil {
			delete(ds.connections, con.Id())
		}
	}()
	con.Close()
}

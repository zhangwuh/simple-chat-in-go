package cs

import "fmt"

type Client interface {
	Connection() Connection
	Close(conn Connection)
}

type DummyClient struct {
	connection Connection
}

func NewDummyClient(conn Connection) *DummyClient{
	defer func() {
		go func() {
			for {
				select {
				case msg, ok :=<- conn.Read():
					if !ok {
						return
					}
					fmt.Println(fmt.Sprintf("[%s:%d]msg from server:%s", conn.Id(), msg.sequence, string(msg.content)))
				}
			}
		}()
	}()
	return &DummyClient{conn}
}

func (dc *DummyClient) Connection() Connection {
	return dc.connection
}

func (dc *DummyClient) Close(conn Connection) {
	conn.Close()
}
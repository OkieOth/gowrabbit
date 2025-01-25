package amqp

import (
	"fmt"
	"sync"

	"github.com/okieoth/gowrabbit/shared/resilence"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ServerOptsFunc func(o *ServerOpts)

type ServerOpts struct {
	Host string
	Port uint
}

func Host(host string) ServerOptsFunc {
	return func(o *ServerOpts) {
		o.Host = host
	}
}

func Port(port uint) ServerOptsFunc {
	return func(o *ServerOpts) {
		o.Port = port
	}
}

func defaultServerOpts() ServerOpts {
	return ServerOpts{
		Host: "localhost",
		Port: 5672,
	}
}

type Server struct {
	ServerOpts
}

func NewServer(fn ...ServerOptsFunc) Server {
	opts := defaultServerOpts()
	for _, f := range fn {
		f(&opts)
	}
	return Server{
		ServerOpts: opts,
	}
}

type ConnectionOptsFunc func(o *ConnectionOpts)

type ConnectionOpts struct {
	User     string
	Password string
	Servers  []Server
}

func User(user string) ConnectionOptsFunc {
	return func(o *ConnectionOpts) {
		o.User = user
	}
}

func Password(pwd string) ConnectionOptsFunc {
	return func(o *ConnectionOpts) {
		o.Password = pwd
	}
}

func Servers(servers []Server) ConnectionOptsFunc {
	return func(o *ConnectionOpts) {
		o.Servers = append(o.Servers, servers...)
	}
}

func defaultConnectionOpts() ConnectionOpts {
	return ConnectionOpts{
		User:     "guest",
		Password: "guest",
		Servers:  make([]Server, 0),
	}
}

type ConnectionState int

const (
	CONNECTED ConnectionState = iota
	DISCONNECTED
)

type Connection struct {
	ConnectionOpts
	mutex      sync.RWMutex
	conn       *amqp.Connection
	connNotify []chan<- ConnectionState
}

func NewConnection(fn ...ConnectionOptsFunc) Connection {
	opts := defaultConnectionOpts()
	for _, f := range fn {
		f(&opts)
	}
	return Connection{
		ConnectionOpts: opts,
	}
}

func (c *Connection) AddConnNotify(n chan<- ConnectionState) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connNotify = append(c.connNotify, n)
}

func (c *Connection) DelConnNotify(n chan<- ConnectionState) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i, v := range c.connNotify {
		if v == n {
			c.connNotify = append(c.connNotify[:i], c.connNotify[i+1:]...)
			break
		}
	}
}

func (c *Connection) ConnNotifyCount() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.connNotify)
}

func (c *Connection) SendConnNotify(s ConnectionState) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for i := 0; i < len(c.connNotify); i++ {
		c.connNotify[i] <- s
	}
}

func (c *Connection) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO - build connection string
	resilentConnectFunc := func() error {
		if conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/"); err == nil {
			fmt.Println("connection established")
			c.conn = conn
			go func() {
				conClosedChan := make(chan *amqp.Error)
				if c.conn != nil {
					c.conn.NotifyClose(conClosedChan)
					fmt.Println("listen for connection clossed events ...")
					if e, ok := <-conClosedChan; ok {
						fmt.Println("connection was closed w/ error: ", e)
					} else {
						fmt.Println("connection closed")
					}
					c.conn = nil
					c.Connect()
				} else {
					// TODO logging
					fmt.Println("connection object is nil")
				}
			}()
			return nil
		} else {
			return err
		}
	}

	if err, tries := resilence.ResilentCall(resilentConnectFunc, 10, 1000, "Rabbitmq-Connect"); err == nil {
		return nil
	} else {
		return fmt.Errorf("finally failed to connect, attempts: %d, reason: %v", tries, err)
	}
}

func (c *Connection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if err := c.conn.Close(); err != nil {
		c.conn = nil
		return nil
	} else {
		return err
	}
}

package amqp

import (
	"fmt"
	"sync"

	"github.com/okieoth/gowrabbit/shared/observer"
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
	connNotify observer.Observer[ConnectionState]
}

func NewConnection(fn ...ConnectionOptsFunc) Connection {
	opts := defaultConnectionOpts()
	for _, f := range fn {
		f(&opts)
	}
	observer := observer.NewObserver[ConnectionState]()
	return Connection{
		ConnectionOpts: opts,
		connNotify:     observer,
	}
}

func (c *Connection) enableReconnect() {
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
			go c.NotifyConnectionClosed()
			c.Connect()
		} else {
			// TODO logging
			fmt.Println("connection object is nil")
		}
	}()
}

func (c *Connection) NotifyConnectionClosed() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connNotify.Notify(DISCONNECTED)
}

func (c *Connection) NotifyConnected() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.connNotify.Notify(CONNECTED)
}

func (c *Connection) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO - build connection string
	resilentConnect := func() error {
		if conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/"); err == nil {
			go c.NotifyConnected()
			fmt.Println("connection established")
			c.conn = conn
			c.enableReconnect()
			return nil
		} else {
			return err
		}
	}

	if err, tries := resilence.ResilentCall(resilentConnect, 10, 1000, "Rabbitmq-Connect"); err == nil {
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

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
	// How many miliseconds should be waited before the next try. With every try the
	// wait time is doubled
	ResilenceWaitMilis int

	// Maximum number of retries in case of connection issues
	ResilenceMaxRetries int
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

func ResilenceMaxRetries(maxRetries int) ConnectionOptsFunc {
	return func(o *ConnectionOpts) {
		o.ResilenceMaxRetries = maxRetries
	}
}

func ResilenceWaitMilis(waitMilis int) ConnectionOptsFunc {
	return func(o *ConnectionOpts) {
		o.ResilenceWaitMilis = waitMilis
	}
}

func defaultConnectionOpts() ConnectionOpts {
	return ConnectionOpts{
		User:                "guest",
		Password:            "guest",
		Servers:             make([]Server, 0),
		ResilenceMaxRetries: 10,
		ResilenceWaitMilis:  1000,
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
	connNotify *observer.Observer[ConnectionState]
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
			if err := c.Connect(); err != nil {
				panic(err)
			}
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

func GetConnectionString(user string, password string, servers []Server) (string, error) {
	if user == "" {
		return "", fmt.Errorf("emtpy user isn't allowed for the connection string")
	}
	if password == "" {
		return "", fmt.Errorf("emtpy password isn't allowed for the connection string")
	}
	if (servers == nil) || (len(servers) == 0) {
		return "", fmt.Errorf("no servers for the connection string given")
	}
	serversStr := ""
	for i, s := range servers {
		if i != 0 {
			serversStr += ","
		}
		serversStr += fmt.Sprintf("%s:%d", s.Host, s.Port)
	}
	// e.g. "amqp://user:password@node1:5672,node2:5672,node3:5672/"
	return fmt.Sprintf("amqp://%s:%s@%s/", user, password, serversStr), nil
}

func (c *Connection) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO - build connection string
	resilentConnect := func() error {
		conStr, err := GetConnectionString(c.User, c.Password, c.Servers)
		if err != nil {
			return fmt.Errorf("error while building connection string: %v", err)
		}
		if conn, err := amqp.Dial(conStr); err == nil {
			go c.NotifyConnected()
			fmt.Println("connection established")
			c.conn = conn
			c.enableReconnect()
			return nil
		} else {
			return err
		}
	}

	if err, tries := resilence.ResilentCall(resilentConnect, c.ResilenceMaxRetries, c.ResilenceWaitMilis, "Rabbitmq-Connect"); err == nil {
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

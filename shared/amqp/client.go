package amqp

import (
	"fmt"
	"sync"

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

type Connection struct {
	ConnectionOpts
	mutex sync.RWMutex
	conn  *amqp.Connection
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

func (c *Connection) Connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// TODO - build connection string
	if conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/"); err == nil {
		c.conn = conn
		return nil
	} else {
		return fmt.Errorf("error, failed to connect to broker: %v", err)
	}
}

// type Client struct {
// 	ClientOpts
// }

// type Connection struct {
// 	m               *sync.Mutex
// 	queueName       string
// 	infolog         *log.Logger
// 	errlog          *log.Logger
// 	connection      *amqp.Connection
// 	channel         *amqp.Channel
// 	done            chan bool
// 	notifyConnClose chan *amqp.Error
// 	notifyChanClose chan *amqp.Error
// 	notifyConfirm   chan amqp.Confirmation
// 	isReady         bool
// }

package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"sync"
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

type Server struct {
	ServerOpts
}

func defaultServerOpts() ServerOpts {
	return ServerOpts{
		Host: "localhost",
		Port: 5672,
	}
}

func NewServer(fn ...ServerOptsFunc) Server {
	ret := defaultServerOpts()
}

type ClientOptsFunc func(o *ClientOpts)

type ClientOpts struct {
	User     string
	Password string
	InfoLog  *log.Logger
	ErrLog   *log.Logger
	Servers  []Server
}

type Client struct {
	ClientOpts
}

type Connection struct {
	m               *sync.Mutex
	queueName       string
	infolog         *log.Logger
	errlog          *log.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

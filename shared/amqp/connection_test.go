package amqp_test

import (
	"testing"

	"github.com/okieoth/gowrabbit/shared/amqp"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s1 := amqp.NewServer()
	assert.Equal(t, "localhost", s1.Host)
	assert.Equal(t, uint(5672), s1.Port)

	s2 := amqp.NewServer(
		amqp.Host("test.com"),
	)
	assert.Equal(t, "test.com", s2.Host)
	assert.Equal(t, uint(5672), s1.Port)

	s3 := amqp.NewServer(
		amqp.Port(8000),
	)
	assert.Equal(t, "localhost", s3.Host)
	assert.Equal(t, uint(8000), s3.Port)

	s4 := amqp.NewServer(
		amqp.Host("test2.com"),
		amqp.Port(8001),
	)
	assert.Equal(t, "test2.com", s4.Host)
	assert.Equal(t, uint(8001), s4.Port)
}

func TestGetConnectionString(t *testing.T) {
	tests := []struct {
		name     string
		user     string
		password string
		servers  []amqp.Server
		want     string
		wantErr  bool
	}{
		// Happy cases
		{
			name:     "Valid inputs with single server",
			user:     "user1",
			password: "pass1",
			servers:  []amqp.Server{amqp.NewServer(amqp.Host("node1"), amqp.Port(5672))},
			want:     "amqp://user1:pass1@node1:5672/",
			wantErr:  false,
		},
		{
			name:     "Valid inputs with multiple servers",
			user:     "user2",
			password: "pass2",
			servers: []amqp.Server{
				amqp.NewServer(amqp.Host("node1"), amqp.Port(5672)),
				amqp.NewServer(amqp.Host("node2"), amqp.Port(5673)),
				amqp.NewServer(amqp.Host("node3"), amqp.Port(5674)),
			},
			want:    "amqp://user2:pass2@node1:5672,node2:5673,node3:5674/",
			wantErr: false,
		},

		// Unhappy cases
		{
			name:     "Empty user",
			user:     "",
			password: "pass1",
			servers:  []amqp.Server{amqp.NewServer(amqp.Host("node1"), amqp.Port(5672))},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Empty password",
			user:     "user1",
			password: "",
			servers:  []amqp.Server{amqp.NewServer(amqp.Host("node1"), amqp.Port(5672))},
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Nil servers slice",
			user:     "user1",
			password: "pass1",
			servers:  nil,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Empty servers slice",
			user:     "user1",
			password: "pass1",
			servers:  []amqp.Server{},
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := amqp.GetConnectionString(tt.user, tt.password, tt.servers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConnectionString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetConnectionString() = %v, want %v", got, tt.want)
			}
		})
	}
}

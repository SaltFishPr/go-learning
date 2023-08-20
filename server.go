package learning

import (
	"bufio"
	"log"
	"net"
	"sync"

	"learning/protocol"
)

type serverOptions struct{}

type ServerOption func(o *serverOptions)

// Server represents an RPC Server.
type Server struct {
	o serverOptions

	serviceMapMu sync.RWMutex
	serviceMap   map[string]*service // map[string]*service

	mu    sync.RWMutex
	conns map[net.Conn]struct{}
}

// NewServer returns a new Server.
func NewServer(opts ...ServerOption) *Server {
	o := serverOptions{}
	for _, opt := range opts {
		opt(&o)
	}

	return &Server{
		o:          o,
		serviceMap: map[string]*service{},
	}
}

// DefaultServer is the default instance of *Server.
var DefaultServer = NewServer()

func (s *Server) ServeConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	m := protocol.NewMessage()
	if err := m.Decode(r); err != nil {
		log.Println(err)
	}
	log.Println(m)
}

// Accept accepts connections on the listener and serves requests
// for each incoming connection. Accept blocks until the listener
// returns a non-nil error. The caller typically invokes Accept in a
// go statement.
func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			return
		}

		s.mu.Lock()
		s.conns[conn] = struct{}{}
		s.mu.Unlock()

		go s.ServeConn(conn)
	}
}

// -----------------------------------------------------------------------------

// Register publishes the receiver's methods in the DefaultServer.
func Register(rcvr any) error {
	return DefaultServer.Register(rcvr)
}

// RegisterName is like Register but uses the provided name for the type
// instead of the receiver's concrete type.
func RegisterName(name string, rcvr any) error {
	return DefaultServer.RegisterName(name, rcvr)
}

// Accept accepts connections on the listener and serves requests
// to DefaultServer for each incoming connection.
// Accept blocks; the caller typically invokes it in a go statement.
func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	Port     int
	listener net.Listener
	handler  Handler
	isClosed atomic.Bool
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isClosed.Load() {
				return
			}
			log.Println("Connection failed")
			continue
		}
		log.Println("Connection established")

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Println(err)
		return
	}

	resWriter := &response.Writer{Conn: conn}
	s.handler(resWriter, req)
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	server := Server{
		Port:     port,
		listener: listener,
		handler:  handler,
	}

	go server.listen()

	return &server, nil
}

package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/debobrad579/httpfromtcp/internal/request"
	"github.com/debobrad579/httpfromtcp/internal/response"
)

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}

func (hErr *HandlerError) Write(w io.Writer) {
	body := []byte(hErr.Message)

	response.WriteStatusLine(w, hErr.StatusCode)
	response.WriteHeaders(w, response.GetDefaultHeaders(len([]byte(hErr.Message))))
	w.Write(body)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

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

	buf := &bytes.Buffer{}
	if hErr := s.handler(buf, req); hErr != nil {
		hErr.Write(conn)
		return
	}

	response.WriteStatusLine(conn, 200)
	response.WriteHeaders(conn, response.GetDefaultHeaders(buf.Len()))
	conn.Write(buf.Bytes())
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

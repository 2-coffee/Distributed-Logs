package web

import (
	"fmt"

	"github.com/2-coffee/Distributed-Logs/server"
	"github.com/valyala/fasthttp"
)

const defaultBufSize = 512 * 1024

type Server struct {
	s *server.InMemory // sending it to our service server
}

// NewServer creates a Server pointer
func NewServer(s *server.InMemory) *Server {
	return &Server{s: s} // creating an instance and sending it to web Server struct
}

func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/write":
		s.writeHandler(ctx)
	case "/read":
		s.readHandler(ctx)
	case "/ack":
		s.ackHandler(ctx)
	default:
		ctx.WriteString("Hello")
	}

}

func (s *Server) writeHandler(ctx *fasthttp.RequestCtx) {
	if err := s.s.Write(ctx.Request.Body()); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // temporary error message; MUST INDICATE THIS LATER
		ctx.WriteString(err.Error())
	}
}

func (s *Server) ackHandler(ctx *fasthttp.RequestCtx) {
	if err := s.s.Ack(); err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // temporary error message; MUST INDICATE THIS LATER
		ctx.WriteString(err.Error())
	}
}

// Read client requests.
func (s *Server) readHandler(ctx *fasthttp.RequestCtx) {
	offset, err := ctx.QueryArgs().GetUint("off") // ask for offset
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("bad `off` GET param: %v", err))
	}

	maxSize, err := ctx.QueryArgs().GetUint("maxSize") // ask for maxSize
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.WriteString(fmt.Sprintf("bad `maxSize` GET param: %v", err))
	}

	err = s.s.Read(uint64(offset), uint64(maxSize), ctx)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError) // temporary error message; MUST INDICATE THIS LATER
		ctx.WriteString(err.Error())
		return
	}

}

func (s *Server) Serve() error {
	return fasthttp.ListenAndServe(":8080", s.handler)
}

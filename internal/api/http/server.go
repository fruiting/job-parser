package http

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	address string
	//handler fasthttp.RequestHandler
	logger *zap.Logger
}

func NewServer(address string /*handler fasthttp.RequestHandler*/, logger *zap.Logger) *Server {
	return &Server{
		address: address,
		//handler: handler,
		logger: logger,
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(s.address, http.HandlerFunc(s.Handler))
}

func (s *Server) Handler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("a")
}

//func (s *Server) ListenAndServe(ctx context.Context) <-chan error {
//	ctx, cancel := context.WithCancel(ctx)
//
//	server := &fasthttp.Server{
//		Handler: s.handler,
//	}
//
//	result := make(chan error)
//	go func() {
//		defer close(result)
//		<-ctx.Done()
//		result <- server.Shutdown()
//	}()
//
//	go func() {
//		err := fasthttp.ListenAndServe(s.address, s.handler)
//		if err != nil {
//			s.logger.Error("can't listenAndServe server", zap.Error(err))
//		}
//		cancel()
//	}()
//
//	return result
//}

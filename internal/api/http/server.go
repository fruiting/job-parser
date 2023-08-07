package http

import (
	"fmt"
	"net/http"
	"net/http/pprof"

	"go.uber.org/zap"
)

const pprofUrlPrefix = "/debug/pprof/"

type Server struct {
	listen      string
	enablePprof bool
	logger      *zap.Logger
}

func NewServer(listen string, enablePprof bool, logger *zap.Logger) *Server {
	return &Server{
		listen:      listen,
		enablePprof: enablePprof,
		logger:      logger,
	}
}

func (s *Server) ListenAndServe() error {
	r := http.NewServeMux()

	r.HandleFunc("/api/v1/ping/", s.ping)

	if s.enablePprof {
		r.HandleFunc(pprofUrlPrefix, pprof.Index)
		r.HandleFunc(fmt.Sprintf("%s/cmdline", pprofUrlPrefix), pprof.Cmdline)
		r.HandleFunc(fmt.Sprintf("%s/profile", pprofUrlPrefix), pprof.Profile)
		r.HandleFunc(fmt.Sprintf("%s/symbol", pprofUrlPrefix), pprof.Symbol)
		r.HandleFunc(fmt.Sprintf("%s/trace", pprofUrlPrefix), pprof.Trace)
	}

	return http.ListenAndServe(s.listen, r)
}

func (s *Server) ping(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PONG"))
}

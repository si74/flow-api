package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Server struct {
	// Has config indicating buffer size
	addr string
	mux  *http.ServeMux
}

// TODO(sneha): add some more config
// TODO(sneha): add logging using logrus
func NewServer(addr string) (*Server, error) {
	mux := http.NewServeMux()

	// Validate addr

	return &Server{
		mux: mux,
	}, nil
}

func Run(ctx context.Context) error {
	// Run multiple goroutines and check if server is cancellable
	return nil
}

func (s *Server) Serve() error {
	s.mux.Handle("/flows", &FlowHandler{})

	// TODO(sneha): add prometheus client endpoint
	// Add go tracing endpoints
	//s.mux.Handle("/metrics", )

	return http.ListenAndServe(s.addr, s.mux)
}

type FlowHandler struct {
	// TODO(sneha): contains data store for flowlist
	// TODO(sneha): add logging
}

func (h *FlowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.handleRead(w, r)
	case "POST":
		h.handleWrite(w, r)
	default:
		// Invalid request type
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *FlowHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Query().Get("hour")
	if str == "" {
		// TODO(sneha): don't want to provide more information to prevent leaky api
		w.WriteHeader(http.StatusBadRequest)
	}
	// Check if hour is a valid int
	hour, err := strconv.Atoi(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	//TODO(sneha): retrieve aggregated information from flow store, marshal to json value, and return
	fmt.Fprintf(w, "successful request for hour: %v", hour)
}

func (h *FlowHandler) handleWrite(w http.ResponseWriter, r *http.Request) {
	// Read body and validate request type

}

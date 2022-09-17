package flowd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/si74/flow-api/internal/store"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	// Has config indicating buffer size
	addr string
	fh   *FlowHandler
}

// TODO(sneha): add some more config
// TODO(sneha): add logging using logrus
// TODO(enable go metrics)
func NewServer(addr string) (*Server, error) {

	// Validate addr

	// Create a flow handler that contains the data structure

	return &Server{
		addr: addr,
		fh:   &FlowHandler{},
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/flows", s.fh)
	// TODO(sneha): add prometheus client endpoint
	// Add go tracing endpoints
	//s.mux.Handle("/metrics", )

	// use this method later with an http custom server and logging middleware
	srv := &http.Server{
		Addr: s.addr,
		// 	Handler: http_logrus.Middleware
		Handler: mux,
		// ErrorLog:
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Handle context cancellation - this is the only way to make the server cancellable
	eg.Go(func() error {
		<-ctx.Done()
		return srv.Close()
	})

	// Separate goroutine to run
	eg.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return eg.Wait()
}

type FlowHandler struct {
	// TODO(sneha): contains data store for flowlist
	// TODO(sneha): add logging
}

func NewFlowHandler() {
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
		return
	}
	// Check if hour is a valid int
	hour, err := strconv.Atoi(str)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//TODO(sneha): retrieve aggregated information from flow store, marshal to json value, and return
	fmt.Fprintf(w, "successful request for hour: %v", hour)
}

// TODO(sneha): Switch from write error to http.Error()
func (h *FlowHandler) handleWrite(w http.ResponseWriter, r *http.Request) {
	// Confirm we are receiving a body type of json
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Read body and validate request type
	var flowList []*store.Flow

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(b, &flowList); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO(sneha): add to flow store database and indicate successful
	fmt.Fprintf(w, "received flows: %v", flowList)
}

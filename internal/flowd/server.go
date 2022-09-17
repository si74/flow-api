package flowd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/si74/flow-api/internal/store"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	// Has config indicating buffer size
	addr string
	fh   *FlowHandler
	ll   *logrus.Logger
	reg  *prometheus.Registry
}

// TODO(sneha): add some more config
// TODO(sneha): add logging using logrus
// TODO(enable go metrics)
func NewServer(addr string, ll *logrus.Logger, reg *prometheus.Registry) (*Server, error) {

	// Validate addr

	// Create a flow handler that contains the data structure

	return &Server{
		addr: addr,
		fh: &FlowHandler{
			ll: ll,
		},
		ll:  ll,
		reg: reg,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/flows", s.fh)
	mux.Handle("/metrics", promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{}))
	// Add go tracing endpoints

	// use this method later with an http custom server and logging middleware
	srv := &http.Server{
		Addr:    s.addr,
		Handler: mux,
		// ErrorLog:
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Handle context cancellation - this is the only way to make the server cancellable
	eg.Go(func() error {
		<-ctx.Done()
		s.ll.Info("context cancellation received")
		return srv.Close()
	})

	// Separate goroutine to run
	eg.Go(func() error {
		s.ll.Infof("starting flowd http server on: %v...", s.addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		s.ll.Info("gracefully stopping flowd server")
		return nil
	})
	return eg.Wait()
}

type FlowHandler struct {
	// TODO(sneha): contains data store for flowlist
	// TODO(sneha): add custom metrics
	ll *logrus.Logger
}

func NewFlowHandler() {
}

func (h *FlowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ll := h.ll.WithField("src", r.RemoteAddr)
	ll.Debug("incoming request")

	switch r.Method {
	case "GET":
		h.handleRead(w, r)
	case "POST":
		h.handleWrite(w, r)
	default:
		// Invalid request type
		h.ll.Debugf("invalid request type %s", r.Method)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *FlowHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	ll := h.ll.WithField("src", r.RemoteAddr)
	ll.Debug("incoming read request")

	str := r.URL.Query().Get("hour")
	if str == "" {
		// TODO(sneha): don't want to provide more information to prevent leaky api
		ll.Debug("read request missing parameter hour")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if hour is a valid int
	hour, err := strconv.Atoi(str)
	if err != nil {
		ll.Debug("read request parameter hour is not an int")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ll.Debug("successful request read request")

	//TODO(sneha): retrieve aggregated information from flow store, marshal to json value, and return
	fmt.Fprintf(w, "successful request for hour: %v", hour)
}

// TODO(sneha): Switch from write error to http.Error()
func (h *FlowHandler) handleWrite(w http.ResponseWriter, r *http.Request) {
	ll := h.ll.WithField("src", r.RemoteAddr)
	ll.Debug("incoming write request")

	// Confirm we are receiving a body type of json
	if r.Header.Get("Content-Type") != "application/json" {
		ll.Debug("invalid write request type: %v", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Read body and validate request type
	var flowList []*store.Flow

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ll.Debug("unble to read write body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(b, &flowList); err != nil {
		ll.Debug("unable to unmarshal write body into flows: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO(sneha): add to flow store database and indicate successful
	fmt.Fprintf(w, "received flows: %v", flowList)
}

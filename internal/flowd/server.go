package flowd

import (
	"context"
	"encoding/json"
	"errors"
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
	// TODO(sneha): Validate addr

	// Create a flow handler that contains the data structure
	fs := store.NewFlowStore(ll)

	return &Server{
		addr: addr,
		fh:   NewFlowHandler(fs, ll),
		ll:   ll,
		reg:  reg,
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
	fs *store.FlowStore
	// TODO(sneha): add custom metrics
	ll *logrus.Logger
}

func NewFlowHandler(fs *store.FlowStore, ll *logrus.Logger) *FlowHandler {
	return &FlowHandler{
		fs: fs,
		ll: ll,
	}
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
		ll.Debugf("read request parameter hour is not an int: %s", str)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ll.Debug("successful request read request")
	flows, err := h.fs.Get(hour)
	if err != nil {
		ll.Debugf("unable to retrieve flows for hour %d: %v", hour, flows)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(flows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ll.Debugf("unable to marshal flows: %v", err)
		return
	}

	w.Write(out)
}

// TODO(sneha): Switch from write error to http.Error()
func (h *FlowHandler) handleWrite(w http.ResponseWriter, r *http.Request) {
	ll := h.ll.WithField("src", r.RemoteAddr)
	ll.Debug("incoming write request")

	// Confirm we are receiving a body type of json
	if r.Header.Get("Content-Type") != "application/json" {
		ll.Debugf("invalid write request type: %v", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
	// Read body and validate request type
	var flowList []*store.Flow

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		ll.Debugf("unble to read write body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(b, &flowList); err != nil {
		ll.Debugf("unable to unmarshal write body into flows: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.fs.Insert(flowList); err != nil {
		ll.Debugf("unable to insert flows: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	grpcClientsBid "github.com/andreyAKor/otus_arch_project/internal/grpc/clients/bid"
)

var (
	ErrServerNotInit  = errors.New("server not init")
	ErrInvalidRequest = errors.New("the request body can’t be pasred as valid data")

	_ io.Closer = (*Server)(nil)
)

type Server struct {
	host      string
	port      int
	bodyLimit int
	bidClient *grpcClientsBid.Bid

	server *http.Server
	ctx    context.Context
}

func New(
	host string,
	port int,
	bodyLimit int,
	bidClient *grpcClientsBid.Bid,
) (*Server, error) {
	return &Server{
		host:      host,
		port:      port,
		bodyLimit: bodyLimit,
		bidClient: bidClient,
	}, nil
}

// Running http-server.
func (s *Server) Run(ctx context.Context) error {
	s.ctx = ctx

	mux := http.NewServeMux()
	mux.HandleFunc("/create", s.method(s.toJSON(s.create), "POST"))
	mux.HandleFunc("/get", s.method(s.toJSON(s.get), "GET"))

	// middlewares
	handler := s.headers(mux)
	handler = s.body(handler)
	handler = s.logger(handler)

	s.server = &http.Server{
		Addr:    net.JoinHostPort(s.host, strconv.Itoa(s.port)),
		Handler: handler,
	}
	if err := s.server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "http-server listen fail")
	}

	return nil
}

func (s *Server) Close() error {
	if s.server == nil {
		//nolint:wrapcheck
		return ErrServerNotInit
	}

	return s.server.Shutdown(s.ctx)
}

// Middleware logger output log info of request, e.g.: r.Method, r.URL etc.
func (s Server) logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newAppResponseWriter(w)

		start := time.Now()
		defer func() {
			i := log.Info()

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				i.Err(err)
			}

			i.Str("ip", host).
				Str("startAt", start.String()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("proto", r.Proto).
				Int("status", rw.statusCode).
				TimeDiff("latency", time.Now(), start)

			if len(r.UserAgent()) > 0 {
				i.Str("userAgent", r.UserAgent())
			}

			i.Msg("http-request")
		}()

		handler.ServeHTTP(rw, r)
	})
}

// Middleware sets http-headers for response.
func (s Server) headers(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS,GET,POST")
		}

		// For OPTIONS requests
		if r.Method == "OPTIONS" {
			return
		}

		// JSON header
		w.Header().Set("Content-Type", "application/json")

		handler.ServeHTTP(w, r)
	})
}

// Middleware preparing body request.
func (s Server) body(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, int64(s.bodyLimit)))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			if err := s.writeJSON(Response{Error: ErrInvalidRequest.Error()}, w); err != nil {
				log.Error().Err(err).Msg("writeJSON fail")
			}

			return
		}

		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		handler.ServeHTTP(w, r)
	})
}

// Checking allowed method for endpoint.
func (s Server) method(handler http.HandlerFunc, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}

		handler(w, r)
	}
}

// Converting Response from endpoint to json-response.
func (s Server) toJSON(h func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rs Response

		data, err := h(w, r)
		if err != nil {
			rs.Error = err.Error()
		} else {
			rs.Data = data
		}

		if err := s.writeJSON(rs, w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Msg("writeJSON fail")

			return
		}
	}
}

// Writing Response structure to json http-response.
func (s Server) writeJSON(rs Response, w io.Writer) error {
	res, err := json.Marshal(&rs)
	if err != nil {
		return errors.Wrap(err, "JSON-marshal fail")
	}

	if _, err := w.Write(res); err != nil {
		return errors.Wrap(err, "write fail")
	}

	return nil
}

var _ http.ResponseWriter = (*appResponseWriter)(nil)

// App wrapper over http.ResponseWriter.
type appResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newAppResponseWriter(w http.ResponseWriter) *appResponseWriter {
	return &appResponseWriter{w, http.StatusOK}
}

func (a *appResponseWriter) WriteHeader(code int) {
	a.statusCode = code
	a.ResponseWriter.WriteHeader(code)
}

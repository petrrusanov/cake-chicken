package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/dimebox/cake-chicken/app/store/service"
	"github.com/dimebox/cake-chicken/app/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

var usernameRegexp = regexp.MustCompile(`<(@[\d\w]+)(\|[\w\d]+)?>`)

// Rest is a rest access server
type Rest struct {
	Version       string
	SharedSecret  string
	DataStore     *service.DataStore

	httpServer *http.Server
	lock       sync.Mutex
}

// Run the lister and request's router, activate rest server
func (s *Rest) Run(httpPort int) {
	log.Printf("[INFO] activate rest HTTP server on port %d", httpPort)

	router := s.routes()

	s.lock.Lock()
	s.httpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", httpPort),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()

	log.Printf("[WARN] http server terminated, %s", err)
}

// Shutdown rest http server
func (s *Rest) Shutdown() {
	log.Print("[WARN] shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	s.lock.Lock()

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[DEBUG] rest shutdown error, %s", err)
		}
	}

	log.Print("[DEBUG] shutdown rest server completed")

	s.lock.Unlock()
}

func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	router.Use(AppInfo("Cake and chicken", s.Version), Ping)

	ipFn := func(ip string) string { return utils.StrongHashValue(ip, s.SharedSecret)[:12] } // logger uses it for anonymization

	router.Route("/", func(rapi chi.Router) {
		rapi.Group(func(ropen chi.Router) {
			ropen.Use(Logger(ipFn, LogBody))
			ropen.Post("/stats", s.stats)
			ropen.Post("/cake/add", s.addCake)
			ropen.Post("/cake/fulfill", s.fulfillCake)
			ropen.Post("/chicken/add", s.addChicken)
			ropen.Post("/chicken/fulfill", s.fulfillChicken)
		})
	})

	return router
}

func encodeJSON(v interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)

	if err := enc.Encode(v); err != nil {
		return nil, errors.Wrap(err, "json encoding failed")
	}

	return buf.Bytes(), nil
}

func renderJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	bytes, err := encodeJSON(v)

	if err != nil {
		return err
	}

	renderJSONFromBytes(w, r, bytes)

	return nil
}

func renderJSONFromBytes(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if status, ok := r.Context().Value(render.StatusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	if _, err := w.Write(data); err != nil {
		log.Printf("[WARN] failed to send response to %s, %s", r.RemoteAddr, err)
	}
}

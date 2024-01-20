package apiserver

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/go-chi/chi/v5"
)

type ApiServerStateCallback struct {
	Started func()
	Stopped func()
}

type ApiServer struct {
	listenaddr string
	router     *chi.Mux
	logger     *log.Logger
	callback   ApiServerStateCallback

	activeThreads *sync.WaitGroup
	server        *http.Server
	state         *atomic.Uint32 // 0 -> Idle, 1 -> Started, 2 -> Stopped

}

func (o *ApiServer) Start() error {
	if o.state.Load() != 0 {
		return errors.New("server already started")
	}

	server := &http.Server{
		Addr:     o.listenaddr,
		Handler:  o.router,
		ErrorLog: o.logger,
	}

	o.server = server

	o.activeThreads.Add(1)
	o.state.Store(1)
	go func() {
		o.callback.Started()
		listenerr := server.ListenAndServe()
		o.logger.Printf("HTTP server ListenAndServe: %v", listenerr)
		o.activeThreads.Done()
	}()

	go func() {
		o.activeThreads.Wait()
		o.server = nil
		o.state.Store(0)
		o.callback.Stopped()
	}()

	return nil
}

func (o *ApiServer) Stop() error {
	if o.state.Load() != 1 {
		return errors.New("server is not in running state")
	}
	o.state.Store(2)
	go func() {
		shutdownerr := o.server.Shutdown(context.Background())
		if shutdownerr != nil {
			o.logger.Printf("Server Shutdown err: %v", shutdownerr)
		}
	}()
	return nil
}

func NewApiServer(listenaddr string, router *chi.Mux, logger *log.Logger, callback ApiServerStateCallback) *ApiServer {
	return &ApiServer{
		listenaddr: listenaddr,
		router:     router,
		logger:     logger,
		callback:   callback,

		activeThreads: &sync.WaitGroup{},
		server:        nil,
		state:         &atomic.Uint32{},
	}
}

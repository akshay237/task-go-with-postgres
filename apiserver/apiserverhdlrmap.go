package apiserver

import (
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type ApiServerHandlerItem struct {
	ApiPathPattern string
	HandlerI       APIHandler
}

type ApiServerHandlerMap struct {
	handlersMap []*ApiServerHandlerItem
	logger      *log.Logger
}

func (o *ApiServerHandlerMap) AddHandler(pattern string, handlerI APIHandler) {
	o.handlersMap = append(o.handlersMap, &ApiServerHandlerItem{
		ApiPathPattern: pattern,
		HandlerI:       handlerI,
	})
}

func (o *ApiServerHandlerMap) ToRouter() *chi.Mux {
	router := chi.NewRouter()

	// 1. Add Logger...
	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  o.logger,
		NoColor: true}))

	// 2. Add panic recoverer...
	router.Use(middleware.Recoverer)

	// 3. Add cors...
	corsI := cors.AllowAll()
	router.Use(corsI.Handler)

	// 4. Create unique requestid for each request...
	router.Use(middleware.RequestID)

	// 5. Add compressor...
	router.Use(middleware.Compress(5, "application/*", "text/*"))

	for _, eachhandlermap := range o.handlersMap {
		router.Route(eachhandlermap.ApiPathPattern, func(r chi.Router) {
			eachhandlermap.HandlerI.RegisterRoutes(r)
		})
	}
	return router
}

func NewApiServerHandlerMap(logger *log.Logger) *ApiServerHandlerMap {
	return &ApiServerHandlerMap{
		handlersMap: make([]*ApiServerHandlerItem, 0, 16),
		logger:      logger,
	}
}

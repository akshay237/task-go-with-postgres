package apiserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func TestApiServer(t *testing.T) {
	apilogger := log.New(os.Stdout, "APISERVER-1", log.Lmicroseconds|log.LstdFlags)
	testlogger := log.New(os.Stdout, "TESTLOGGER1", log.Lmicroseconds|log.LstdFlags)

	listenaddr := ":11000"
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.AllowContentType("application/json"))

	middlewarelogger := middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger:  log.New(os.Stdout, "APISERVER-1", log.Lmicroseconds|log.LstdFlags),
		NoColor: false})
	router.Use(middlewarelogger)
	router.Use(middleware.Recoverer)

	router.Get("/error", func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			apilogger.Printf("Header Key:%v, Value:%v", name, values)
		}
		APIResponseBadRequest(w, r, "BAD_REQUEST", "Random data", "Sending random error")
	})

	router.Get("/key/user1", func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			apilogger.Printf("Header Key:%v, Value:%v", name, values)
		}
		// Read the full body...
		bodycontent, readerr := io.ReadAll(r.Body)
		var mapdata map[string]interface{}
		err := json.Unmarshal(bodycontent, &mapdata)
		apilogger.Printf("Unmarshal error: Err:%v", err)
		apilogger.Printf("Unrmarshalled data:Body:%v", mapdata)
		w.Header().Set("Content-Type", "application/json")
		if readerr != nil {
			w.WriteHeader(400)
			w.Write([]byte(`{"key":"value"}`))
		}
		w.WriteHeader(200)
		w.Write(bodycontent)
		w.Write([]byte("Hello world"))
		w.Write([]byte("next row"))
		time.Sleep(time.Millisecond * 1000)
		w.Write([]byte("After timeout"))
	})

	activeThreads := &sync.WaitGroup{}

	apiservercallback := ApiServerStateCallback{
		Started: func() {
			testlogger.Println("Api server started")
		},
		Stopped: func() {
			testlogger.Println("Api server stopped")
			activeThreads.Done()
		},
	}

	apiserver := NewApiServer(listenaddr, router, apilogger, apiservercallback)
	activeThreads.Add(1)
	starterr := apiserver.Start()
	if starterr != nil {
		t.Fatal("Api server start failed", starterr)
	}

	time.Sleep(time.Millisecond * 60000)
	testlogger.Println("Issuing Stop")
	apiserver.Stop()
	activeThreads.Wait()
	testlogger.Println("Test stopped")
}

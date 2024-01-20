package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"task-go-with-postgres/apiserver"
	"task-go-with-postgres/config"
	"task-go-with-postgres/database"
	"time"

	"github.com/natefinch/lumberjack"
	zlog "github.com/rs/zerolog/log"
)

func getServerConfig() (*config.ServerConfig, error) {
	var configFileName string
	flag.StringVar(&configFileName, "configfile", "", "Please pass the config file")

	flag.Parse()
	if configFileName == "" {
		log.Fatal("Please pass the config file as the argument")
	}
	zlog.Info().Msg("Config File to be used:" + configFileName)
	return config.ParseConfigFromFile(configFileName)
}

func setupLogger(logDirPath string) {
	logFilePath := path.Join(logDirPath, "general.log")
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(logFilePath),
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	log.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile)

	log.SetOutput(lumberjackLogger)
}

func addLoggerFile(logDir string, logFile string, logger *log.Logger) {
	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(path.Join(logDir, logFile)),
		MaxSize:    5, // MB
		MaxBackups: 5,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}
	logger.SetFlags(log.Ltime | log.Ldate | log.Lmicroseconds | log.Lshortfile)
	logger.SetOutput(lumberjackLogger)
}

func waitTillStopFile(stoppedflag *uint32, stopch chan string, stopfilepath string) {
	log.Println("Stopfile Path:", stopfilepath)
	for atomic.LoadUint32(stoppedflag) == 0 {
		if _, err := os.Stat(stopfilepath); err == nil {
			log.Println("Removing stopfile")
			err := os.Remove(stopfilepath)
			if err != nil {
				log.Println("Removing stopfile has failed")
			}
			break
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
	}
	stopch <- "stopfilereceived"
}

func main() {
	runtime.GOMAXPROCS(2)
	serverConfig, err := getServerConfig()
	if err != nil {
		log.Fatal("Server Config Error", err)
	}

	logDirPath := path.Join(*serverConfig.ScratchDir, "logs")
	setupLogger(logDirPath)

	postgresDBI := database.NewPostgresDBService(serverConfig.PostgresConfig)
	fmt.Println("PostgresDB", postgresDBI)

	apilogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "apiserver.log", apilogger)

	servlogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "service.log", servlogger)

	hdlrlogger := log.New(os.Stdout, "", log.Lmicroseconds|log.LstdFlags|log.Llongfile)
	addLoggerFile(logDirPath, "handler.log", hdlrlogger)

	// 1. Service

	// 2. Handlers

	activeThreads := &sync.WaitGroup{}

	apiservercallback := apiserver.ApiServerStateCallback{
		Started: func() {
			apilogger.Println("Api server started")
		},
		Stopped: func() {
			apilogger.Println("Api server stopped")
			activeThreads.Done()
		},
	}

	listenaddr := fmt.Sprintf(":%d", serverConfig.APIServerConfig.Port)
	handlerMap := apiserver.NewApiServerHandlerMap(apilogger)
	handlerMap.AddHandler("/api/v1/task", nil)

	apiserver := apiserver.NewApiServer(listenaddr, handlerMap.ToRouter(), apilogger, apiservercallback)

	activeThreads.Add(1)
	starterr := apiserver.Start()
	if starterr != nil {
		apilogger.Fatal("Api server start failed", starterr)
	}

	ossigch := make(chan os.Signal, 1)
	signal.Notify(ossigch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	stopfilech := make(chan string, 1)
	stoppedflag := uint32(0)
	pendingthreads := &sync.WaitGroup{}
	pendingthreads.Add(1)
	go func() {
		waitTillStopFile(&stoppedflag, stopfilech, path.Join(*serverConfig.ScratchDir, "stopfile"))
		atomic.StoreUint32(&stoppedflag, 1)
		pendingthreads.Done()
	}()

	select {
	case signalrcvd := <-ossigch:
		log.Println("Signal received from OS:", signalrcvd.String())
	case <-stopfilech:
		log.Println("Stop file seen, so starting the stop task")
	}
	atomic.StoreUint32(&stoppedflag, 1)

	stoperr := apiserver.Stop()
	if stoperr != nil {
		apilogger.Fatal("Api server stop failed", stoperr)
	}

	pendingthreads.Wait()
	activeThreads.Wait()
	log.Println("Api Server shutdown successfully...")

}

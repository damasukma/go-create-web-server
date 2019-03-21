package main

import (
	"context"
	_ "errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "reflect"
	_ "strings"
	"sync/atomic"
	"time"
)

var (
	listen string
	healty int32
)

func main() {

	flag.StringVar(&listen, "port", "8000", "serve listen address")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	logger.Println("Server Is Starting...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})

	serve := fmt.Sprintf(":%s", listen)

	ServeResult := &http.Server{
		Addr:         serve,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		logger.Println("Server is shutting down..")
		atomic.StoreInt32(&healty, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := ServeResult.Shutdown(ctx); err != nil {
			logger.Fatal("Could not gracefully shutdown the server: %v\n", err)
		}

		close(done)

	}()

	//--------------------Serve Start---------------------------//

	log.Printf("Listen To Port %s", listen)
	atomic.StoreInt32(&healty, 1)

	if err := ServeResult.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Cannot Connect To Serve Port %s", serve)
	}

	<-done
	logger.Println("Server Stoped")

}

// func IsFlagDefine(name string) bool {
// 	check := false

// 	flag.Visit(func(f *flag.Flag) {
// 		if f.Name == name {
// 			check = true
// 		}
// 	})

// 	return check
// }

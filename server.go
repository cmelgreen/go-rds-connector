
package main

import (
	"context"
	"log"
	"os"
	"time"
	"github.com/julienschmidt/httprouter"
)

const (
	heartbeatTime = 10
)

var (
	logOut = os.Stdout
	logPrefix = log.Prefix()
	logFlags = log.Flags()
)

// Server struct for storing database, mux, and logger
type Server struct{
    db *Database
    mux *httprouter.Router
    log *log.Logger
}

func newServer(ctx context.Context) *Server {
	s := Server{
        log: log.New(logOut, logPrefix, logFlags),
        mux: httprouter.New(),
    }

    db, err := ConnectToDB(ctx)
    if err != nil {
        s.log.Println(err)
    }

	s.db = db

	s.maintainDBConnection(ctx)

	return &s
}

func (s *Server) maintainDBConnection(ctx context.Context) {
	go func() {
		var err error
		for {
			if s.db.Connected(ctx) != true {
				s.db, err = ConnectToDB(ctx)
				if err != nil {
					s.log.Println(err)
				}
			}

			time.Sleep(heartbeatTime * time.Second)
		}
	}()
}
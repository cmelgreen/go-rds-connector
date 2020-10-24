package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Index is a closure that returns a function ro check the database connection and write status to user
func (s *Server) index() httprouter.Handle {
    return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
        if s.db.Connected(r.Context()) {
			fmt.Fprint(w, "Connected!")
        } else {
            http.Error(w, "Error connecting to database", http.StatusNotFound)
            s.log.Println("Error connectiong to database")
        }
    }
}
package main

import (
	"fmt"
	"net/http"
	"time"
)

func ServeStatic(addr string, hashes *map[string]string, dir string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for path, hash := range *hashes {
			fmt.Fprintf(w, "%s %s\n", path, hash)
		}
	})

	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir(dir))))

	s := &http.Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		Handler:        nil,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()

	if err != nil {
		return err
	}

	return nil
}

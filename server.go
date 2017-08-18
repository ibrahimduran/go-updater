package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MainHandler struct {
	secretKey      string
	fileServer     http.Handler
	checksumServer http.HandlerFunc
}

func (h MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.secretKey != "" && h.secretKey != r.URL.Query().Get("secret") {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	switch {
	case strings.HasPrefix(r.URL.Path, "/data"):
		h.fileServer.ServeHTTP(w, r)
	case r.URL.Path == "/":
		h.checksumServer(w, r)
	default:
		http.NotFound(w, r)
	}
}

func ServeStatic(addr string, secretKey string, hashes *map[string]string, dir string) error {
	fileServer := http.StripPrefix("/data", http.FileServer(http.Dir(dir)))

	checksumServer := func(w http.ResponseWriter, r *http.Request) {
		for path, hash := range *hashes {
			fmt.Fprintf(w, "%s %s\n", path, hash)
		}
	}

	http.Handle("/", MainHandler{secretKey, fileServer, checksumServer})

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

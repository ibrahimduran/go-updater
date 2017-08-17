package main

import (
	"log"
	"net/http"
)

func ServeStatic(addr string, dir string) {
	err := http.ListenAndServe(addr, http.FileServer(http.Dir(dir)))
	if err != nil {
		log.Fatal(err)
	}
}

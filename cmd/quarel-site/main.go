// Command quarel-site serves the prebuilt static Astro site.
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	root := flag.String("root", "/site", "static site directory")
	flag.Parse()
	files := http.FileServer(http.Dir(*root))
	http.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	http.Handle("/", files)
	log.Printf("serving site on %s from %s", *addr, *root)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Printf("site server: %v", err)
		os.Exit(1)
	}
}

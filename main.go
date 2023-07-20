package main

import (
	"log"
	"net/http"

	gddo "github.com/golang/gddo/httputil"
)

func handleHTML(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving HTML at %s", req.URL.Path)
	w.Header().Set("Content-Type", "text/html")
	w.Header().Add("Link", `</>; rel="preload"; as="fetch"; type="application/json"; crossorigin="anonymous"`)
	w.Header().Add("Link", `</data.json>; rel="preload"; as="fetch"; type="application/json"; crossorigin="anonymous"`)
	http.ServeFile(w, req, "main.html")
}

func handleJSON(w http.ResponseWriter, req *http.Request) {
	log.Printf("Serving JSON at %s", req.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, req, "main.json")
}

func handle(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" && req.URL.Path != "/data.json" {
		http.NotFound(w, req)
		return
	}

	w.Header().Add("Vary", "Accept")
	contentType := gddo.NegotiateContentType(req, []string{"text/html", "application/json"}, "text/html")
	switch contentType {
	case "text/html":
		handleHTML(w, req)
	case "application/json":
		handleJSON(w, req)
	}
}

func main() {
	log.Print("Listening on :8000")
	err := http.ListenAndServe(":8000", http.HandlerFunc(handle))
	log.Printf("%s", err)
}

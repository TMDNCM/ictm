package main

import (
	"github.com/TMDNCM/ictm/service/web"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", web.NewHandler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

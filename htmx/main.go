package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/kokweikhong/go-playground/htmx/router"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {

	r := router.InitRouter()

	slog.Info("Server started at port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}

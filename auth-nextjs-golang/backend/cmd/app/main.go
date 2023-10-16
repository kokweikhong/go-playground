package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/kokweikhong/go-playground/auth-nextjs-golang/backend/internal/router"
)

func main() {
	fmt.Println("Hello World")
	log.Default().SetFlags(log.Lshortfile | log.LstdFlags)

	r := router.Init()

	slog.Info("Starting server on port 8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}

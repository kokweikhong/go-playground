package main

import (
	"net/http"

	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/config"
	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/model"
	"github.com/kokweikhong/go-playground/chi-middleware-jwt-auth/router"
)

func main() {

	if err := config.InitConfig(".env"); err != nil {
		panic(err)
	}

	err := model.InitUsers()

	if err != nil {
		panic(err)
	}
	r := router.InitRouter()

	http.ListenAndServe(":8080", r)
}

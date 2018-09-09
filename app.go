package main

import (
	"log"
	"net/http"

	"github.com/TV4/graceful"
	calculator "github.com/ichsanrp/tax-calculator/src/taxCalculator"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var err error
	router := httprouter.New()

	_, err = calculator.Init(router)
	if err != nil {
		log.Fatal(err)
		return
	}

	router.ServeFiles("/*filepath", http.Dir("public"))

	graceful.LogListenAndServe(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
}

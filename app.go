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

	// serving index.html to public. this file will be used for UI example
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "public/index.html")
	})

	// we store our calculator in module and passing http router to it's controller
	_, err = calculator.Init(router)
	if err != nil {
		log.Fatal(err)
		return
	}

	// gracefull shutdown to finish all request before shutdown
	graceful.LogListenAndServe(&http.Server{
		Addr:    ":8080",
		Handler: router,
	})
}

package main

import (
	"sqlrest/internal/config"
	"sqlrest/internal/handler"
	"sqlrest/internal/util"
	"flag"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {

	configFile := flag.String("f", "config.yaml", "Path to the config file")
	flag.Parse()

	config.LoadConfig(*configFile)

	util.LogInit()

	util.ConnectDB()
	defer util.MainDB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api", handler.ApiHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/help", handler.HelpHandler).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./dist/index.html")
	})
	router.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("./dist/css"))))
	router.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("./dist/js"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

}

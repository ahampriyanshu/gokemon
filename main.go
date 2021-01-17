package main

import (
	"fmt"
	"math/rand"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"github.com/mtslzr/pokeapi-go"
	"github.com/joho/godotenv"
)
	
	var tpl = template.Must(template.ParseFiles("index.html"))
	
	func indexHandler(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
	
	func searchHandler(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		params := u.Query()
		searchQuery := params.Get("q")
		if searchQuery == "" {
			searchQuery =strconv.Itoa(rand.Intn(500 - 1) + 1) 
		}
		pokedata, err := pokeapi.Pokemon(searchQuery)
		fmt.Println("Search Query is: ", pokedata)


	}

	func main() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
		}
	
		port := os.Getenv("PORT")
		if port == "" {
			port = "3000"
		}

		fs := http.FileServer(http.Dir("assets"))
		mux := http.NewServeMux()
		mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
		mux.HandleFunc("/", indexHandler)
		mux.HandleFunc("/search", searchHandler)
		http.ListenAndServe(":"+port, mux)
	}
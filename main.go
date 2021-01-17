package main

import (
	"bytes"
	"math/rand"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"strconv"
	"github.com/mtslzr/pokeapi-go"
	"github.com/joho/godotenv"
)
	
	var tpl = template.Must(template.ParseFiles("index.html"))

	type Pokedata struct {
		Query      string
		Name      string
		Experience  int
		Height      int
		Weight      int
		Avatar      string
		Type		string
		// Abilities   []struct
	}
	
	func dataHandler(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		params := u.Query()
		searchQuery := params.Get("q")
		searchQuery = strings.ToLower(searchQuery) 
		if searchQuery == "" {
			searchQuery = strconv.Itoa(rand.Intn(500 - 1) + 1) 
		}

		pokejson, err := pokeapi.Pokemon(searchQuery)

		pokedata := &Pokedata{
			Query:        searchQuery,
			Name:    	  strings.Title(pokejson.Name),
			Experience:   pokejson.BaseExperience,
			Height:       pokejson.Height,
			Weight:       pokejson.Weight,
			Avatar:  	  pokejson.Sprites.FrontDefault,
			Type:		  pokejson.Types[0].Type.Name,
			// Abilities:	  pokejson.Abilities[0],
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, pokedata)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}

	buf.WriteTo(w)
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
		mux.HandleFunc("/", dataHandler)
		mux.HandleFunc("/search", dataHandler)
		http.ListenAndServe(":"+port, mux)
	}
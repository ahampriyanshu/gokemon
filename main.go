package main

import (
	"bytes"
	"math/rand"
	"html/template"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"strconv"
	"github.com/mtslzr/pokeapi-go"
	"github.com/joho/godotenv"
	"time"
)
	
	var tpl = template.Must(template.ParseFiles("index.html"))

	type Pokedata struct {
		Query       string
		Name        string
		Experience  int
		Height      int
		Weight      int
		Avatar      string
		Type		string
	}
	
	func dataHandler(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rand.Seed(time.Now().UTC().UnixNano())
		params := u.Query()
		searchQuery := params.Get("q")
		searchQuery = strings.ToLower(searchQuery) 
		if searchQuery == "" {
			searchQuery = strconv.Itoa(rand.Intn(500 - 1) + 1) 
		}

		pokejson, err := pokeapi.Pokemon(searchQuery)

		if err != nil {
		pokejson,err = pokeapi.Pokemon( strconv.Itoa(rand.Intn(500 - 1) + 1))
		}


		pokedata := &Pokedata{
			Query:        searchQuery,
			Name:    	  strings.Title(pokejson.Name),
			Experience:   pokejson.BaseExperience,
			Height:       pokejson.Height,
			Weight:       pokejson.Weight,
			Avatar:  	  pokejson.Sprites.FrontDefault,
			Type:		  pokejson.Types[0].Type.Name,
		}

		buf := &bytes.Buffer{}
		err = tpl.Execute(buf, pokedata)
		if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
		}

	buf.WriteTo(w)
	}

	func sendSW(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("sw.js")
		if err != nil {
			http.Error(w, "Couldn't read file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		w.Write(data)
	}
	
	func sendManifest(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile("manifest.json")
		if err != nil {
			http.Error(w, "Couldn't read file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(data)
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
		mux.HandleFunc("/sw.js", sendSW)
		mux.HandleFunc("/manifest.json", sendManifest)
		mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
		mux.HandleFunc("/", dataHandler)
		mux.HandleFunc("/search", dataHandler)
		http.ListenAndServe(":"+port, mux)
	}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type Person struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Teacher struct {
	Person
	Degree string
}

type Envelope map[string]any

func deleteFromSlice(slice *[]any, indexToDelete int) {
	if slice != nil && indexToDelete >= 0 && indexToDelete < len(*slice) {
		*slice = append((*slice)[:indexToDelete], (*slice)[indexToDelete+1:]...)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var persons = []Person{}
	persons = append(persons,
		Person{1, "Ryan", 20},
		Person{2, "Riley", 40},
		Person{3, "Alice", 60},
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("Hello World!")
	})

	r.Get("/persons", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(persons)
	})

	r.Get("/persons/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idToFind, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		for i := 0; i < len(persons); i++ {
			curr := persons[i]
			if int64(curr.Id) == idToFind {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(curr)
				return
			}
		}

		w.WriteHeader(http.StatusNotFound)
	})

	r.Post("/persons", func(w http.ResponseWriter, r *http.Request) {
		var j Person
		err := json.NewDecoder(r.Body).Decode(&j)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("There was an issue with your message"))
			return
		}

		j.Id = len(persons) + 1
		persons = append(persons, j)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(persons)
	})

	r.Delete("/persons/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		personID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			w.WriteHeader(500)
		}

		indexToDelete := -1
		for i, p := range persons {
			if int64(p.Id) == personID {
				indexToDelete = i
				break
			}
		}

		if indexToDelete == -1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		persons = append(persons[:indexToDelete], persons[indexToDelete+1:]...)
		w.WriteHeader(http.StatusOK)
	})

	fmt.Printf("Server is now running on port: %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}

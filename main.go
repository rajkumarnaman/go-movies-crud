package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func swagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.ServeFile(w, r, "swagger.json")
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	//converting json to its own format
	w.Header().Set("Content-Type", "application/json")
	//New Encoder Function-->gets a JSON encoding of any type and encodes/writes it any writable stream that implements a io
	json.NewEncoder(w).Encode(movies) //Encode is used to encode a value to json format

}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//variables can be retrieved using mux.vars
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(100000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	//params
	//loop over the movies
	//delete the movie with the id that is sent
	//add a new movie - the movie that we send in the body of postman
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = strconv.Itoa(rand.Intn(100000000))
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", index).Methods("GET")
	router.HandleFunc("/swagger.json", swagger).Methods("GET")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))

	r := mux.NewRouter() //used to create a router instance

	movies = append(movies, Movie{ID: "1", Isbn: "831004", Title: "Movie One", Director: &Director{Firstname: "Anurag", Lastname: "Kashyap"}})
	movies = append(movies, Movie{ID: "2", Isbn: "831005", Title: "Movie Two", Director: &Director{Firstname: "Rohit", Lastname: "Shetty"}})
	r.HandleFunc("/movies", getMovies).Methods("GET") //use the HandleFunc method of your router instance to assign routes to handler functions along with the request type that the handler function handles
	//"/movies-->route""getMovies-->function name""GET-->method"
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting Server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))

} //makefile,docker,cubernate,aws,etc

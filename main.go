package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Cars represents database entity.
type Car struct {
	ID    string `json:"id"`
	Model string `json:"model"`
	Age   string `json:"age"`
	Price string `json:"price"`
}

const port = ":8080"

var err error

// Check for errors
func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var db *sql.DB

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "WELCOME")
}

// Returns a list of all database cars to the response.
func GetCars(w http.ResponseWriter, r *http.Request) {
	rows, e := db.Query("SELECT * FROM cars")
	checkErr(e)
	defer rows.Close()
	cars := make([]*Car, 0)

	for rows.Next() {
		c := new(Car)
		e := rows.Scan(&c.ID, &c.Model, &c.Age, &c.Price)
		checkErr(e)
		cars = append(cars, c)
		if e = rows.Err(); e != nil {
			log.Fatal(e)
		}
	}
	json.NewEncoder(w).Encode(cars)
}

// Returns a single database car matching given ID parameter.
func GetCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(500), 500)
	}
	row := db.QueryRow("SELECT * FROM cars WHERE ID = ?", id)
	e := row.Scan(&car.ID, &car.Model, &car.Age, &car.Price)
	checkErr(e)
	json.NewEncoder(w).Encode(car)
}

// Create car into the database.
func CreateCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	_ = json.NewDecoder(r.Body).Decode(&car)
	_, e := db.Exec("INSERT INTO cars (model, age, price) VALUES (?, ?, ?)", car.Model, car.Age, car.Price)
	checkErr(e)
	json.NewEncoder(w).Encode(car)
}

// Update car(identified by parameter) from the database.
func UpdateCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(500), 500)
	}
	row := db.QueryRow("SELECT * FROM cars WHERE ID = ?", id)
	e := row.Scan(&car.ID, &car.Model, &car.Age, &car.Price)
	checkErr(e)
	_ = json.NewDecoder(r.Body).Decode(&car)
	_, err := db.Exec("UPDATE cars SET model=?, age=?, price=? WHERE ID = ?", car.Model, car.Age, car.Price, id)
	checkErr(err)
	json.NewEncoder(w).Encode(car)
}

// Removes car (identified by parameter) from the database.
func DeleteCar(w http.ResponseWriter, r *http.Request) {
	var car Car
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, http.StatusText(500), 500)
	}
	_, e := db.Exec("DELETE FROM cars WHERE ID = ?", id)
	checkErr(e)
	json.NewEncoder(w).Encode(car)
}

func main() {
	fmt.Println("first line works")
	// Establish a connection to MySQL.
	db, err = sql.Open("mysql", "root:root@tcp(db:3306)/cars")
	checkErr(err)

	defer db.Close()

	err = db.Ping()
	checkErr(err)

	// Create routes
	r := mux.NewRouter()

	r.HandleFunc("/", Index)
	r.HandleFunc("/cars", GetCars).Methods("GET")
	r.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	r.HandleFunc("/cars/{id}", CreateCar).Methods("POST")
	r.HandleFunc("/cars/{id}", UpdateCar).Methods("PUT")
	r.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	log.Println("Server is up on " + port + " port")
	log.Fatal(http.ListenAndServe(port, r))

}

package main

import (
  "fmt"
  "log"
  "net/http"
  _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
  "database/sql"
  "encoding/json"
)

// Person struct
type Users struct {
	Id        string `form:"id" json:"id"`
	FirstName string `form:"firstname" json:"firstname"`
	LastName  string `form:"lastname" json:"lastname"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Users
}

// function untuk connect ke database
func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:admin3kali@tcp(localhost:3306)/golang")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

// funtion untuk memparsing data MySQL ke JSON
func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	var users Users //mapping variable user
	var arr_user []Users
	var response Response

	db := connect()
	defer db.Close()

	rows, err := db.Query("select id,first_name,last_name from person")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&users.Id, &users.FirstName, &users.LastName); err != nil {
			log.Fatal(err.Error())

		} else {
			arr_user = append(arr_user, users)
		}
	}

	response.Status = 1
	response.Message = "Success"
	response.Data = arr_user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func insertUsersMultipart(w http.ResponseWriter, r *http.Request) {

  var response Response

	db := connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")

	_, err = db.Exec("INSERT INTO person (first_name, last_name) values (?,?)",
		first_name,
		last_name,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Add"
	log.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func updateUsersMultipart(w http.ResponseWriter, r *http.Request) {

  var response Response

	db := connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id := r.FormValue("user_id")
	first_name := r.FormValue("first_name")
	last_name := r.FormValue("last_name")

	_, err = db.Exec("UPDATE person set first_name = ?, last_name = ? where id = ?",
		first_name,
		last_name,
		id,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Update Data"
	log.Print("Update data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func deleteUsersMultipart(w http.ResponseWriter, r *http.Request) {

  var response Response

	db := connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	id := r.FormValue("user_id")

	_, err = db.Exec("DELETE from person where id = ?",
		id,
	)

	if err != nil {
		log.Print(err)
	}

	response.Status = 1
	response.Message = "Success Delete Data"
	log.Print("Delete data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}


func main() {

  router := mux.NewRouter()
  router.HandleFunc("/users",returnAllUsers).Methods("GET")
  router.HandleFunc("/users",insertUsersMultipart).Methods("POST")
  router.HandleFunc("/users",updateUsersMultipart).Methods("PUT")
  router.HandleFunc("/users",deleteUsersMultipart).Methods("DELETE")
  http.Handle("/", router)
  fmt.Println("Connected to port 1234")
  log.Fatal(http.ListenAndServe(":1234", router))
}

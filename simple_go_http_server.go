package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func writeJSON(w http.ResponseWriter, value any) {
	b, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)
	w.Write(b)
	w.Write([]byte("\n"))
}

func readJSON[T any](r *http.Request, out *T) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var parsed_request T
	err = json.Unmarshal(b, &parsed_request)
	if err != nil {
		return err
	}
	*out = parsed_request
	return nil
}

type listResponseMessage struct {
	Values []int `json:"values"`
}

var list_stmt *sql.Stmt

func listValues(w http.ResponseWriter, r *http.Request) {
	var rows *sql.Rows
	var err error
	rows, err = list_stmt.Query()
	if err != nil {
		w.WriteHeader(500)
		return
	}
	var response listResponseMessage
	for rows.Next() {
		var value int
		rows.Scan(&value)
		response.Values = append(response.Values, value)
	}
	writeJSON(w, response)
}

type pushRequestMessage struct {
	Value int `json:"value"`
}

var push_stmt *sql.Stmt

func pushValue(w http.ResponseWriter, r *http.Request) {
	var err error
	var request pushRequestMessage
	err = readJSON[pushRequestMessage](r, &request)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("parsing error\n"))
		return
	}
	_, err = push_stmt.Exec(request.Value)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
}

func check_starting_error(err error) {
	if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "sqlite3.db")
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS sample_table ( id INTEGER PRIMARY KEY AUTOINCREMENT, value INTEGER NOT NULL);")
	check_starting_error(err)
	list_stmt, err = db.Prepare("SELECT (value) FROM sample_table;")
	check_starting_error(err)
	push_stmt, err = db.Prepare("INSERT INTO sample_table (value) VALUES ($1);")
	check_starting_error(err)

	http.HandleFunc("/list", listValues)
	http.HandleFunc("/push", pushValue)

	err = http.ListenAndServe(":8080", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else {
		check_starting_error(err)
	}
}

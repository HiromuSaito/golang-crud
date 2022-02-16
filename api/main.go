package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Println("application start!")
	//テスト用
	http.HandleFunc("/test", handleAllTestRequest)
	http.HandleFunc("/users", handleAllUserRequest)
	http.Handle("/users/", http.StripPrefix("/users/", http.HandlerFunc(handleSingleUserRequest)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

//テスト用
func handleAllTestRequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	w.WriteHeader(http.StatusOK)
}

func handleAllUserRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		postUser(w, r)
	default:
		http.Error(w, r.Method+" method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query(`select * from users order by id;`)

	if err != nil {
		log.Printf("select user err:%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var users []User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.Name); err != nil {
			log.Fatal("scan err:", err)
		}
		users = append(users, user)
	}

	body, _ := json.MarshalIndent(&users, "", "  ")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func postUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	if _, err := Db.Exec(`insert into users (
	id
 ,name
)
select
 case
    when max(id) is null then 1
    else max(id)+1
  end
 ,?
from
    users;`, user.Name); err != nil {
		log.Printf("user insert err:%s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleSingleUserRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPut:
		putUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, r.Method+" method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	var user User

	if err := Db.QueryRow("select * from users where id = ?", id).Scan(&user.Id, &user.Name); err != nil {
		log.Printf("user find err:%s,id=%s\n", err, id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, _ := json.MarshalIndent(&user, "", "  ")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func putUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	if _, err := Db.Exec("update users set name = ? where id = ?", user.Name, id); err != nil {
		log.Printf("user update err:%s,id=%s\n", err, id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path
	if _, err := Db.Exec("delete from users where id = ?", id); err != nil {
		log.Printf("user delete err:%s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type User struct {
	Id   int    `db:"id",json:"id"`
	Name string `db:"name",json:"name"`
}

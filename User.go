package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	IdUser    int
	NamaUser  string
	Alamat    string
	Pekerjaan string
}

func DbConnect() (db *sql.DB) {
	DbDriver := "mysql"
	DbUser := "root"
	DbHost := "tcp(localhost)"
	DbName := "webappgolang"
	db, err := sql.Open(DbDriver, DbUser+":"+"@"+DbHost+"/"+DbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var templ = template.Must(template.ParseGlob("user/*"))

func Index(w http.ResponseWriter, r *http.Request) {

	db := DbConnect()
	selDB, err := db.Query("SELECT * FROM user ORDER BY NamaUser DESC")
	if err != nil {
		panic(err.Error())
	}
	usr := User{}
	res := []User{}
	for selDB.Next() {
		var IdUser int
		var NamaUser, Alamat, Pekerjaan string

		err = selDB.Scan(&IdUser, &NamaUser, &Alamat, &Pekerjaan)
		if err != nil {
			panic(err.Error())
		}
		usr.IdUser = IdUser
		usr.NamaUser = NamaUser
		usr.Alamat = Alamat
		usr.Pekerjaan = Pekerjaan
		res = append(res, usr)
	}
	templ.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}
func Input(w http.ResponseWriter, r *http.Request) {
	templ.ExecuteTemplate(w, "Form", nil)
}
func Insert(w http.ResponseWriter, r *http.Request) {
	db := DbConnect()
	log.Println(r.Method)
	if r.Method == "POST" {
		UserName := r.FormValue("UserName")
		Alamat := r.FormValue("Alamat")
		Pekerjaan := r.FormValue("Pekerjaan")
		inserData, err := db.Prepare("INSERT INTO user(NamaUser,Alamat,Pekerjaan) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		inserData.Exec(UserName, Alamat, Pekerjaan)
		log.Println("INSERT Sucess ")
	}
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}
func Edit(w http.ResponseWriter, r *http.Request) {

	db := DbConnect()
	id := r.URL.Query().Get("IdUser")
	selDb, err := db.Query("SELECT * FROM user WHERE IdUser=?", id)
	if err != nil {
		panic(err.Error())
	}

	usr := User{}
	for selDb.Next() {
		var IdUser int
		var NamaUser, Alamat, Pekerjaan string
		err = selDb.Scan(&IdUser, &NamaUser, &Alamat, &Pekerjaan)
		if err != nil {
			panic(err.Error())
		}
		usr.IdUser = IdUser
		usr.NamaUser = NamaUser
		usr.Alamat = Alamat
		usr.Pekerjaan = Pekerjaan

	}
	templ.ExecuteTemplate(w, "Edit", usr)
	defer db.Close()
}

func Update(w http.ResponseWriter, r *http.Request) {
    db := DbConnect()
    if r.Method == "POST" {
		UserId := r.FormValue("IdUser")
        UserName := r.FormValue("UserName")
		Alamat := r.FormValue("Alamat")
		Pekerjaan := r.FormValue("Pekerjaan")
        inserData, err := db.Prepare("UPDATE User SET NamaUser=?, Alamat=?, Pekerjaan=? WHERE IdUser=?")
        if err != nil {
            panic(err.Error())
        }
        inserData.Exec(UserName, Alamat, Pekerjaan, UserId)
        log.Println("UPDATE: Name: " + UserName + " | Alamat: " + Alamat)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
    db := DbConnect()
    emp := r.URL.Query().Get("IdUser")
    delForm, err := db.Prepare("DELETE FROM User WHERE IdUser=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}


func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/Input", Input)
	http.HandleFunc("/Insert", Insert)
	http.HandleFunc("/Edit", Edit)
	http.HandleFunc("/Update",Update)
	http.HandleFunc("/Delete",Delete)
	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}

package service

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/wedancedalot/squirrel"
	"PentagoServer/db"
)

const UsersTable = "users"

type User struct {
	Id	uint64 `db:"usr_id"`
	Email string `db:"usr_email"`
	Password string `db:"usr_password"`
}

func CreateUser(writer http.ResponseWriter, router *http.Request) {
	decoder := json.NewDecoder(router.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte("422 - Cannot read the request's body."))
		return
	}

	data := squirrel.Eq{
		"usr_email":           user.Email,
		"usr_password":   user.Password,
	}

	_,err = db.GetClientMysql().Insert(squirrel.Insert(UsersTable).SetMap(data))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte("500 - Cannot create user."))
		return
	}
	js, err := json.Marshal(map[string]interface{}{
		"status": true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return

}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user User
	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Cannot read the request's body."))
		return
	}

	var selectedUser User
	q:= squirrel.Select("*").From(UsersTable).Where(squirrel.Eq{"usr_email": user.Email})

	err = db.GetClientMysql().FindFirst(&selectedUser,q)
	if err == db.ErrNoRows{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - User not found."))
		return
	}
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Cannot create user."))
		return
	}

	if user.Password != selectedUser.Password{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Wrong password."))
		return
	}

	js, err := json.Marshal(map[string]interface{}{
		"status": true,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func Results(writer http.ResponseWriter, router *http.Request) {
	//todo
}

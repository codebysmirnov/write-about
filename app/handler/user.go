package handler

import (
	"encoding/json"
	"factory/in_progress/app/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// getUserOr404 check user in db
func getUserOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, id).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

type createUserIn struct {
	Login    string `json:login`
	Password string `json:password`
}

// CreateUser - create new user on system
func CreateUser(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	in := &createUserIn{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.TrimSpace(in.Login)) < 2 {
		respondError(w, http.StatusBadRequest, "Login must have more charachters")
		return
	}

	if len(strings.TrimSpace(in.Login)) < 6 {
		respondError(w, http.StatusBadRequest, "Password must be longer than five charachters")
		return
	}

	newUser := &model.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := db.Create(newUser).Error; err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	responseJSON(w, http.StatusCreated, newUser)
}

// AllUsers - get all users on system
func AllUsers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	db.Find(&users)
	responseJSON(w, http.StatusOK, users)
}

// GetUserByID from table users
func GetUserByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	User := getUserOr404(db, id, w, r)
	if User == nil {
		return
	}
	responseJSON(w, http.StatusOK, User)
}

// CheckApii - get all users on system
func CheckApii(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	println("PING")
	users := []model.User{}
	db.Find(&users)
	responseJSON(w, http.StatusOK, users)
}

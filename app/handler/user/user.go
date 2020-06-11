package user

import (
	"encoding/json"
	"factory/in_progress/app/handler"
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
		handler.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

type createUserIn struct {
	Login    string `json:login`
	Password string `json:password`
}

// User struct
type User struct {
	db *gorm.DB
}

// New user module
func New(db *gorm.DB) *User {
	return &User{db: db}
}

// CreateUser - create new user on system
func (u *User) CreateUser(w http.ResponseWriter, r *http.Request) {
	in := &createUserIn{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.TrimSpace(in.Login)) < 2 {
		handler.RespondError(w, http.StatusBadRequest, "Login must have more charachters")
		return
	}

	if len(strings.TrimSpace(in.Login)) < 6 {
		handler.RespondError(w, http.StatusBadRequest, "Password must be longer than five charachters")
		return
	}

	newUser := &model.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := u.db.Create(newUser).Error; err != nil {
		handler.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	handler.ResponseJSON(w, http.StatusCreated, newUser)
}

// AllUsers - get all users on system
func (u *User) AllUsers(w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	u.db.Find(&users)
	handler.ResponseJSON(w, http.StatusOK, users)
}

// GetUserByID from table users
func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	User := getUserOr404(u.db, id, w, r)
	if User == nil {
		return
	}
	handler.ResponseJSON(w, http.StatusOK, User)
}

// Register handler
func (u *User) Register(r *mux.Router) {
	r.HandleFunc("/", u.CreateUser).Methods("POST")
	r.HandleFunc("/", u.AllUsers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", u.GetUserByID).Methods("GET")
}

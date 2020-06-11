package user

import (
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

// getUserOr404 check user in db
// TODO: Delete this method, and use ctx for get user
func getUserOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *model.User {
	user := model.User{}
	if err := db.First(&user, id).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &user
}

// User struct
type User struct {
	db *gorm.DB
}

// New user module
func New(db *gorm.DB) *User {
	return &User{db: db}
}

// AllUsers - get all users on system
func (u *User) AllUsers(w http.ResponseWriter, r *http.Request) {
	users := []model.User{}
	u.db.Find(&users)
	utils.ResponseJSON(w, http.StatusOK, users)
}

// GetUserByID from table users
func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, _ := strconv.Atoi(vars["id"])

	User := getUserOr404(u.db, id, w, r)
	if User == nil {
		return
	}
	utils.ResponseJSON(w, http.StatusOK, User)
}

// Register route handlers
func (u *User) Register(r *mux.Router) {
	r.HandleFunc("/", u.AllUsers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", u.GetUserByID).Methods("GET")
}

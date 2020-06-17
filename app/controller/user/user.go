package user

import (
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
)

// User struct
type User struct {
	db *gorm.DB
}

// New user module
func New(db *gorm.DB) *User {
	if db == nil {
		panic("failed to initialize Auth controller: db parameter is nil-pointer")
	}
	return &User{db: db}
}

// AllUsers - get all users on system
func (u *User) AllUsers(w http.ResponseWriter, r *http.Request) {
	var users []model.User
	u.db.Find(&users)
	utils.ResponseJSON(w, http.StatusOK, users)
}

// GetUserByID from table users
func (u *User) GetUserByID(w http.ResponseWriter, r *http.Request) {
	var userMeta = r.Context().Value("user").(auth.Meta)

	user := model.User{}
	if err := u.db.First(&user, userMeta["user_id"].(uint)).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusOK, user)
}

// Register controller handlers
func (u *User) Register(r *mux.Router) {
	r.HandleFunc("/", u.AllUsers).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", u.GetUserByID).Methods("GET")
}

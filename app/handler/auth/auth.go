package auth

import (
	"encoding/json"
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"net/http"
	"strings"
)

// Auth struct
type Auth struct {
	db   *gorm.DB
	auth auth.Auth
}

// New Auth module
func New(db *gorm.DB, auth auth.Auth) *Auth {
	return &Auth{db: db, auth: auth}
}

type userCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Registration - create new user on this great system
func (a *Auth) Registration(w http.ResponseWriter, r *http.Request) {
	in := &userCredentials{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.TrimSpace(in.Login)) < 2 {
		utils.RespondError(w, http.StatusBadRequest, "Login must have more characters")
		return
	}

	if len(strings.TrimSpace(in.Login)) < 6 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be longer than five characters")
		return
	}

	newUser := &model.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := a.db.Create(newUser).Error; err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, newUser)
}

// Login - take a token
func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	in := &userCredentials{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(strings.TrimSpace(in.Login)) < 2 {
		utils.RespondError(w, http.StatusBadRequest, "Login must have more characters")
		return
	}

	if len(strings.TrimSpace(in.Login)) < 6 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be longer than five characters")
		return
	}

	user := &model.User{}
	if err := a.db.Model(user).Where("login = ? and password = ?", in.Login, in.Password).First(user).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	uc := map[string]interface{}{"user_id": user.ID}
	token, err := a.auth.Generate(uc)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, token)
}

// Register route handlers
func (a *Auth) Register(r *mux.Router) {
	r.HandleFunc("/register", a.Registration).Methods("POST")
	r.HandleFunc("/login", a.Login).Methods("POST")
}
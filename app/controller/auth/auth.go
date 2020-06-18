package auth

import (
	"database/sql"
	"encoding/json"
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"github.com/gorilla/mux"
	"github.com/volatiletech/sqlboiler/boil"
	"net/http"
	"strings"
)

// Auth struct
type Auth struct {
	db   *sql.DB
	auth auth.Auth
}

// New Auth controller
// may throw panic when one of parameter is nil-pointer
func New(db *sql.DB, auth auth.Auth) *Auth {
	if db == nil {
		panic("failed to initialize Auth controller: db parameter is nil-pointer")
	}
	if auth == nil {
		panic("failed to initialize Auth controller: auth parameter is nil-pointer")
	}
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

	if len(strings.TrimSpace(in.Password)) < 6 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be longer than five characters")
		return
	}

	newUser := &model.User{
		Login:    in.Login,
		Password: in.Password,
	}

	if err := newUser.Insert(a.db, boil.Infer()); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, utils.Empty{})
}

type loginResponse struct {
	Token string `json:"token"`
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

	if len(strings.TrimSpace(in.Password)) < 6 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be longer than five characters")
		return
	}

	user, err := model.Users(
		model.UserWhere.Login.EQ(in.Login),
		model.UserWhere.Password.EQ(in.Password),
	).One(a.db)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	uc := auth.Meta{"user_id": user.ID}
	token, err := a.auth.Generate(uc)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, &loginResponse{Token: token})
}

// Register controller handlers
func (a *Auth) Register(r *mux.Router) {
	r.HandleFunc("/register", a.Registration).Methods("POST")
	r.HandleFunc("/login", a.Login).Methods("POST")
}

package diary

import (
	"encoding/json"
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Diary struct
type Diary struct {
	db *gorm.DB
}

// New diary module
func New(db *gorm.DB) *Diary {
	if db == nil {
		panic("failed to initialize Auth controller: db parameter is nil-pointer")
	}
	return &Diary{db: db}
}

// CreateDiary - create new diary for current user
func (d *Diary) CreateDiary(w http.ResponseWriter, r *http.Request) {
	var userMeta = r.Context().Value("user").(auth.Meta)

	userId := userMeta["user_id"].(uint)

	user := model.User{}
	if err := d.db.First(&user, userId).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	diary := model.Diary{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&diary); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := d.db.Save(&diary).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusCreated, nil)
}

// GetDiary user get only his diary
func (d *Diary) GetDiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var userMeta = r.Context().Value("user").(auth.Meta)
	year, _ := strconv.Atoi(vars["year"])

	diary := model.Diary{}
	if err := d.db.First(&diary, userMeta["user_id"].(uint), year).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusOK, diary)
}

// Register controller handlers
func (d *Diary) Register(r *mux.Router) {
	r.HandleFunc("/", d.CreateDiary).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", d.GetDiary).Methods("GET")
}

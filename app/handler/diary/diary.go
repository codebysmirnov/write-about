package diary

import (
	"encoding/json"
	"factory/in_progress/app/model"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// getUserOr404 check diary in db
func getDiaryOr404(db *gorm.DB, userID int, year int, w http.ResponseWriter, r *http.Request) *model.Diary {
	diary := model.Diary{}
	if err := db.First(&diary, userID, year).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &diary
}

// CreateDiary - create new diary for current user
func CreateDiary(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	id := r.Form.Get("id_user")
	userID, err := strconv.Atoi(id)
	if err != nil {
		return
	}

	user := getUserOr404(db, userID, w, r)
	if user == nil {
		return
	}

	diary := model.Diary{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&diary); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&diary).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	responseJSON(w, http.StatusCreated, diary)
}

// GetDiary user get only his diary
// TODO: Write token system, for auth user
func GetDiary(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, _ := strconv.Atoi(vars["id"])
	year, _ := strconv.Atoi(vars["year"])

	User := getDiaryOr404(db, userID, year, w, r)
	if User == nil {
		return
	}
	responseJSON(w, http.StatusOK, User)
}

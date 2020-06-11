package diary

import (
	"factory/in_progress/app/model"
	"factory/in_progress/app/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// getUserOr404 check diary in db
func getDiaryOr404(db *gorm.DB, userID int, year int, w http.ResponseWriter, r *http.Request) *model.Diary {
	diary := model.Diary{}
	if err := db.First(&diary, userID, year).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &diary
}

// Diary struct
type Diary struct {
	db *gorm.DB
}

// New diary module
func New(db *gorm.DB) *Diary {
	return &Diary{db: db}
}

// CreateDiary - create new diary for current user
func (d *Diary) CreateDiary(w http.ResponseWriter, r *http.Request) {
	//r.ParseForm()
	//
	//id := r.Form.Get("id_user")
	//userID, err := strconv.Atoi(id)
	//if err != nil {
	//	return
	//}

	//user := getUserOr404(d.db, userID, w, r)
	//if user == nil {
	//	return
	//}
	//
	//diary := model.Diary{}
	//decoder := json.NewDecoder(r.Body)
	//
	//if err := decoder.Decode(&diary); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, err.Error())
	//	return
	//}
	//defer r.Body.Close()
	//
	//if err := d.db.Save(&diary).Error; err != nil {
	//	utils.RespondError(w, http.StatusInternalServerError, err.Error())
	//	return
	//}
	utils.ResponseJSON(w, http.StatusCreated, nil)
}

// GetDiary user get only his diary
// TODO: Write token system, for auth user
func (d *Diary) GetDiary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userID, _ := strconv.Atoi(vars["id"])
	year, _ := strconv.Atoi(vars["year"])

	User := getDiaryOr404(d.db, userID, year, w, r)
	if User == nil {
		return
	}
	utils.ResponseJSON(w, http.StatusOK, User)
}

// Register route handlers
func (d *Diary) Register(r *mux.Router) {
	r.HandleFunc("/", d.CreateDiary).Methods("POST")
	r.HandleFunc("/{id:[0-9]+}", d.GetDiary).Methods("GET")
}

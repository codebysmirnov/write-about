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

type createDiary struct {
	Year int `json:"year"`
}

// CreateDiary - create new diary for current user
func (d *Diary) CreateDiary(w http.ResponseWriter, r *http.Request) {
	in := &createDiary{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if in.Year <= 0 {
		utils.RespondError(w, http.StatusBadRequest, "Invalid year number")
		return
	}

	var userMeta = r.Context().Value("user").(auth.Meta)

	userId := userMeta["user_id"]

	user := model.User{}
	if err := d.db.Where("id = ?", userId).First(&user).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	newDiary := &model.Diary{
		Year:   in.Year,
		UserID: user.ID,
	}

	if err := d.db.Save(&newDiary).Error; err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, nil)
}

type ListItem struct {
	ID   uint `json:"id"`
	Year int  `json:"year"`
}

type getDiaryResponse struct {
	Diaries []ListItem `json:"diaries"`
}

// GetDiary user get only his diary by year
func (d *Diary) GetDiary(w http.ResponseWriter, r *http.Request) {
	var userMeta = r.Context().Value("user").(auth.Meta)
	year, _ := strconv.Atoi(r.FormValue("year"))

	query := d.db.Where("user_id = ?", userMeta["user_id"])
	if year > 0 {
		query = query.Where("year = ?", year)
	}

	var diaries []model.Diary

	if err := query.Find(&diaries).Error; err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	var items []ListItem
	for _, diary := range diaries {
		items = append(items, ListItem{diary.ID, diary.Year})
	}

	utils.ResponseJSON(w, http.StatusOK, getDiaryResponse{items})
}

// Register controller handlers
func (d *Diary) Register(r *mux.Router) {
	r.HandleFunc("/", d.CreateDiary).Methods("POST")
	r.HandleFunc("/search", d.GetDiary).Methods("GET")
}

package diary

import (
	"database/sql"
	"encoding/json"
	"github.com/codebysmirnov/write-about/app/middleware/auth"
	"github.com/codebysmirnov/write-about/app/model"
	"github.com/codebysmirnov/write-about/app/utils"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Diary struct
type Diary struct {
	db *sql.DB
}

// New diary module
func New(db *sql.DB) *Diary {
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

	user, err := model.Users(qm.Where("id = ?", userId)).One(d.db)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	newDiary := &model.Diary{
		Year:   in.Year,
		UserID: null.IntFrom(user.ID),
	}

	if err := newDiary.Insert(d.db, boil.Infer()); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, nil)
}

type ListItem struct {
	ID   int `json:"id"`
	Year int `json:"year"`
}

type getDiaryResponse struct {
	Diaries []ListItem `json:"diaries"`
}

// GetDiary user get only his diary by year
func (d *Diary) GetDiary(w http.ResponseWriter, r *http.Request) {
	var userMeta = r.Context().Value("user").(auth.Meta)
	year, _ := strconv.Atoi(r.FormValue("year"))

	query := []qm.QueryMod{qm.Where("user_id = ?", userMeta["user_id"])}
	if year > 0 {
		query = append(query, model.DiaryWhere.Year.EQ(year))
	}

	diaries, err := model.Diaries(query...).All(d.db)
	if err != nil {
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

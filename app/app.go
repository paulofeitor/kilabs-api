package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/paulofeitor/kilabs-api/app/routes"
	"github.com/paulofeitor/kilabs-api/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

type App struct {
	Router *httprouter.Router
	DB     *sql.DB
}

func (a *App) Initialize(config *config.Config) {
	a.setDatabase(config)
	a.Router = httprouter.New()
	a.setRoutes()
}

func (a *App) setRoutes() {
	a.Router.GET("/candidate", a.GetAllCandidates)
	a.Router.POST("/candidate", a.AddCandidate)
	a.Router.GET("/candidate/:candidate_id", a.GetCandidate)
	a.Router.PUT("/candidate/:candidate_id", a.UpdateCandidate)
	a.Router.DELETE("/candidate/:candidate_id", a.DeleteCandidate)

	a.Router.GET("/candidate/:candidate_id/slot", a.GetCandidateSlots)
	a.Router.POST("/candidate/:candidate_id/slot", a.AddCandidateSlot)
	a.Router.PUT("/candidate/:candidate_id/slot/:slot_id", a.UpdateCandidateSlot)
	a.Router.DELETE("/candidate/:candidate_id/slot/:slot_id", a.DeleteCandidateSlot)

	a.Router.GET("/interviewer", a.GetAllInterviewers)
	a.Router.POST("/interviewer", a.AddInterviewer)
	a.Router.GET("/interviewer/:interviewer_id", a.GetInterviewer)
	a.Router.PUT("/interviewer/:interviewer_id", a.UpdateInterviewer)
	a.Router.DELETE("/interviewer/:interviewer_id", a.DeleteInterviewer)

	a.Router.GET("/interviewer/:interviewer_id/slot", a.GetInterviewerSlots)
	a.Router.POST("/interviewer/:interviewer_id/slot", a.AddInterviewerSlot)
	a.Router.PUT("/interviewer/:interviewer_id/slot/:slot_id", a.UpdateInterviewerSlot)
	a.Router.DELETE("/interviewer/:interviewer_id/slot/:slot_id", a.DeleteInterviewerSlot)

	a.Router.POST("/slot", a.SlotMatching)
}

/* CANDIDATES */
func (a *App) AddCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.AddCandidate(a.DB, w, r, ps)
}
func (a *App) GetCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetCandidate(a.DB, w, r, ps)
}
func (a *App) GetAllCandidates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetAllCandidates(a.DB, w, r, ps)
}
func (a *App) UpdateCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.UpdateCandidate(a.DB, w, r, ps)
}
func (a *App) DeleteCandidate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.DeleteCandidate(a.DB, w, r, ps)
}

/* CANDIDATES */
/* CANDIDATES SLOTS */
func (a *App) AddCandidateSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.AddCandidateSlot(a.DB, w, r, ps)
}
func (a *App) GetCandidateSlots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetCandidateSlots(a.DB, w, r, ps)
}
func (a *App) UpdateCandidateSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.UpdateCandidateSlot(a.DB, w, r, ps)
}
func (a *App) DeleteCandidateSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.DeleteCandidateSlot(a.DB, w, r, ps)
}

/* CANDIDATES SLOTS */
/* INTERVIEWERS */
func (a *App) AddInterviewer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.AddInterviewer(a.DB, w, r, ps)
}
func (a *App) GetInterviewer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetCandidate(a.DB, w, r, ps)
}
func (a *App) GetAllInterviewers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetAllInterviewers(a.DB, w, r, ps)
}
func (a *App) UpdateInterviewer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.UpdateInterviewer(a.DB, w, r, ps)
}
func (a *App) DeleteInterviewer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.DeleteInterviewer(a.DB, w, r, ps)
}

/* INTERVIEWERS */
/* INTERVIEWERS SLOTS */
func (a *App) AddInterviewerSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.AddInterviewerSlot(a.DB, w, r, ps)
}
func (a *App) GetInterviewerSlots(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.GetInterviewerSlots(a.DB, w, r, ps)
}
func (a *App) UpdateInterviewerSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.UpdateInterviewerSlot(a.DB, w, r, ps)
}
func (a *App) DeleteInterviewerSlot(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.DeleteInterviewerSlot(a.DB, w, r, ps)
}

/* INTERVIEWERS SLOTS */
/* SLOT MATCH */
func (a *App) SlotMatching(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	routes.SlotMatching(a.DB, w, r, ps)
}

/* SLOT MATCH */

// Setting database
func (a *App) setDatabase(config *config.Config) {
	var err error
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	a.DB, err = sql.Open(config.DB.Driver, dbDSN)
	if err != nil {
		log.Fatal("Could not connect database")
	}
	err = a.DB.Ping()
	if err != nil {
		log.Fatal("Could not ping database")
	}
}

func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

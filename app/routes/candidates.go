package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/paulofeitor/kilabs-api/app/model"
)

func AddCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	candidate := model.Candidate{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&candidate); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	query := "INSERT INTO candidates VALUES (NULL, ?, NOW());"
	result, err := db.Exec(query, candidate.Name)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	candidateId, err := result.LastInsertId()
	if err != nil {
		log.Println("Database Last Insert Id Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	candidate.Id = int(candidateId)

	writeJSON(w, http.StatusOK, candidate)
	return
}

func GetAllCandidates(db *sql.DB, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	candidates := []model.Candidate{}
	query := "SELECT id, name FROM candidates;"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer rows.Close()
	for rows.Next() {
		candidate := model.Candidate{}
		err = rows.Scan(&candidate.Id, &candidate.Name)
		if err != nil {
			log.Println("Database Scan Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		candidates = append(candidates, candidate)
	}
	if len(candidates) == 0 {
		writeError(w, http.StatusNoContent, "No Content")
		return
	}
	writeJSON(w, http.StatusOK, candidates)
	return
}

func GetCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	candidate := model.Candidate{}
	query := "SELECT id, name FROM candidates WHERE id = ?;"
	err := db.QueryRow(query, ps.ByName("candidate_id")).Scan(&candidate.Id, &candidate.Name)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, candidate)
	return
}

func UpdateCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	candidate := model.Candidate{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&candidate); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	candidate.Id, _ = strconv.Atoi(ps.ByName("candidate_id"))
	query := "UPDATE candidates SET name = ? WHERE id = ?;"
	_, err := db.Exec(query, candidate.Name, candidate.Id)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, candidate)
	return
}

func DeleteCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := "DELETE FROM candidates WHERE id = ?;"
	_, err := db.Exec(query, ps.ByName("candidate_id"))
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, nil)
	return
}

func AddCandidateSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slot := model.Slot{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slot); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	slot.PersonId, _ = strconv.Atoi(ps.ByName("candidate_id"))

	query := "INSERT INTO slots SET candidate_id = ?, initial_time = ?, final_time = ?"
	result, err := db.Exec(query, slot.PersonId, slot.InitialTime, slot.FinalTime)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	slotId, err := result.LastInsertId()
	if err != nil {
		log.Println("Database Last Insert Id Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	for _, weekday := range slot.Weekdays {
		query = "INSERT INTO slots_weekdays SET slot_id = ?, weekday = ?;"
		_, err = db.Exec(query, slotId, weekday)
		if err != nil {
			log.Println("Database Query Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	slot.Id = int(slotId)
	writeJSON(w, http.StatusOK, slot)
	return
}

func GetCandidateSlots(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slots := []model.Slot{}
	candidateId := ps.ByName("candidate_id")
	query := "SELECT id, initial_time, final_time FROM slots WHERE candidate_id = ?"
	rows, err := db.Query(query, candidateId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer rows.Close()

	for rows.Next() {
		slot := model.Slot{}
		slot.PersonId, _ = strconv.Atoi(candidateId)
		err = rows.Scan(&slot.Id, &slot.InitialTime, &slot.FinalTime)
		if err != nil {
			log.Println("Database Scan Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		weekdays := []time.Weekday{}
		query = "SELECT weekday FROM slots_weekdays WHERE slot_id = ?"
		rowsWeekdays, err := db.Query(query, slot.Id)
		if err != nil {
			log.Println("Database Query Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		defer rowsWeekdays.Close()
		for rowsWeekdays.Next() {
			var weekday time.Weekday
			err = rowsWeekdays.Scan(&weekday)
			if err != nil {
				log.Println("Database Scan Error ::", err.Error())
				writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			weekdays = append(weekdays, weekday)
		}
		slot.Weekdays = weekdays
		slots = append(slots, slot)
	}

	writeJSON(w, http.StatusOK, slots)
}

func UpdateCandidateSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slot := model.Slot{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slot); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	slot.Id, _ = strconv.Atoi(ps.ByName("slot_id"))
	slot.PersonId, _ = strconv.Atoi(ps.ByName("candidate_id"))

	query := "UPDATE slots SET initial_time = ?, final_time = ? WHERE id = ?"
	_, err := db.Exec(query, slot.InitialTime, slot.FinalTime, slot.Id)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	query = "DELETE FROM slots_weekdays WHERE slot_id = ?"
	_, err = db.Exec(query, slot.Id)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	for _, weekday := range slot.Weekdays {
		query = "INSERT INTO slots_weekdays SET slot_id = ?, weekday = ?;"
		_, err = db.Exec(query, slot.Id, weekday)
		if err != nil {
			log.Println("Database Query Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
	}

	writeJSON(w, http.StatusOK, slot)
	return
}

func DeleteCandidateSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slotId := ps.ByName("slot_id")
	query := "DELETE FROM slots_weekdays WHERE slot_id = ?"
	_, err := db.Exec(query, slotId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	query = "DELETE FROM slots WHERE id = ?"
	_, err = db.Exec(query, slotId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeJSON(w, http.StatusOK, nil)
	return
}

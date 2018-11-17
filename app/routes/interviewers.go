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

func AddInterviewer(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	interviewer := model.Interviewer{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&interviewer); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	query := "INSERT INTO interviewers VALUES (NULL, ?, NOW());"
	result, err := db.Exec(query, interviewer.Name)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	interviewerId, err := result.LastInsertId()
	if err != nil {
		log.Println("Database Last Insert Id Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	interviewer.Id = int(interviewerId)

	writeJSON(w, http.StatusOK, interviewer)
	return
}

func GetAllInterviewers(db *sql.DB, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	interviewers := []model.Interviewer{}
	query := "SELECT id, name FROM interviewers;"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer rows.Close()
	for rows.Next() {
		interviewer := model.Interviewer{}
		err = rows.Scan(&interviewer.Id, &interviewer.Name)
		if err != nil {
			log.Println("Database Scan Error ::", err.Error())
			writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		interviewers = append(interviewers, interviewer)
	}
	if len(interviewers) == 0 {
		writeError(w, http.StatusNoContent, "No Content")
		return
	}
	writeJSON(w, http.StatusOK, interviewers)
	return
}

func GetInterviewer(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	interviewer := model.Interviewer{}
	query := "SELECT id, name FROM interviewers WHERE id = ?;"
	err := db.QueryRow(query, ps.ByName("interviewer_id")).Scan(&interviewer.Id, &interviewer.Name)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, interviewer)
	return
}

func UpdateInterviewer(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	interviewer := model.Interviewer{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&interviewer); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	interviewer.Id, _ = strconv.Atoi(ps.ByName("interviewer_id"))
	query := "UPDATE interviewers SET name = ? WHERE id = ?;"
	_, err := db.Exec(query, interviewer.Name, interviewer.Id)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, interviewer)
	return
}

func DeleteInterviewer(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	query := "DELETE FROM interviewers WHERE id = ?;"
	_, err := db.Exec(query, ps.ByName("interviewer_id"))
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	writeJSON(w, http.StatusOK, nil)
	return
}

func AddInterviewerSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slot := model.Slot{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slot); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	slot.PersonId, _ = strconv.Atoi(ps.ByName("interviewer_id"))

	query := "INSERT INTO slots SET interviewer_id = ?, initial_time = ?, final_time = ?"
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

func GetInterviewerSlots(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slots := []model.Slot{}
	interviewerId := ps.ByName("interviewer_id")
	query := "SELECT id, initial_time, final_time FROM slots WHERE interviewer_id = ?"
	rows, err := db.Query(query, interviewerId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	defer rows.Close()

	for rows.Next() {
		slot := model.Slot{}
		slot.PersonId, _ = strconv.Atoi(interviewerId)
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

func UpdateInterviewerSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slot := model.Slot{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&slot); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	slot.Id, _ = strconv.Atoi(ps.ByName("slot_id"))
	slot.PersonId, _ = strconv.Atoi(ps.ByName("interviewer_id"))

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

func DeleteInterviewerSlot(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

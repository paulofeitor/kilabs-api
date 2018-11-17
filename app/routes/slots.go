package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/paulofeitor/kilabs-api/app/model"
)

func SlotMatching(db *sql.DB, w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	request := model.SlotMatchingRequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		log.Println("Bad Request")
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	defer r.Body.Close()

	candidateSlots, err := getCandidateSlots(db, request.Candidate.Id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	validSlots := []model.Slot{}
	for _, candidateSlot := range candidateSlots {
		validSlotForCandidateSlot := model.Slot{}
		for _, interviewer := range request.Interviewers {
			interviewerSlots, err := getInterviewerSlots(db, interviewer.Id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			interviewerSlotFound := false
			for _, interviewerSlot := range interviewerSlots {
				valid := false
				validSlot := model.Slot{}
				if validSlotForCandidateSlot.InitialTime == "" {
					valid, validSlot = match(candidateSlot, interviewerSlot)
				} else {
					valid, validSlot = match(validSlotForCandidateSlot, interviewerSlot)
				}
				if valid {
					interviewerSlotFound = true
					validSlotForCandidateSlot = validSlot
				}
			}
			if !interviewerSlotFound {
				validSlotForCandidateSlot = model.Slot{}
				continue
			}
		}
		if validSlotForCandidateSlot.InitialTime != "" {
			validSlots = append(validSlots, validSlotForCandidateSlot)
		}
	}
	validSlots = splitSlots(validSlots)
	writeJSON(w, http.StatusOK, validSlots)
}

func getCandidateSlots(db *sql.DB, candidateId int) ([]model.Slot, error) {
	slots := []model.Slot{}
	query := "SELECT id, initial_time, final_time FROM slots WHERE candidate_id = ?"
	rows, err := db.Query(query, candidateId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		return slots, err
	}
	defer rows.Close()

	for rows.Next() {
		slot := model.Slot{}
		err = rows.Scan(&slot.Id, &slot.InitialTime, &slot.FinalTime)
		if err != nil {
			log.Println("Database Scan Error ::", err.Error())
			return slots, err
		}

		weekdays := []time.Weekday{}
		query = "SELECT weekday FROM slots_weekdays WHERE slot_id = ?"
		rowsWeekdays, err := db.Query(query, slot.Id)
		if err != nil {
			log.Println("Database Query Error ::", err.Error())
			return slots, err
		}
		defer rowsWeekdays.Close()
		for rowsWeekdays.Next() {
			var weekday time.Weekday
			err = rowsWeekdays.Scan(&weekday)
			if err != nil {
				log.Println("Database Scan Error ::", err.Error())
				return slots, err
			}
			weekdays = append(weekdays, weekday)
		}
		slot.Weekdays = weekdays
		slots = append(slots, slot)
	}
	return slots, nil
}

func getInterviewerSlots(db *sql.DB, interviewerId int) ([]model.Slot, error) {
	slots := []model.Slot{}
	query := "SELECT id, initial_time, final_time FROM slots WHERE interviewer_id = ?"
	rows, err := db.Query(query, interviewerId)
	if err != nil {
		log.Println("Database Query Error ::", err.Error())
		return slots, err
	}
	defer rows.Close()

	for rows.Next() {
		slot := model.Slot{}
		err = rows.Scan(&slot.Id, &slot.InitialTime, &slot.FinalTime)
		if err != nil {
			log.Println("Database Scan Error ::", err.Error())
			return slots, err
		}

		weekdays := []time.Weekday{}
		query = "SELECT weekday FROM slots_weekdays WHERE slot_id = ?"
		rowsWeekdays, err := db.Query(query, slot.Id)
		if err != nil {
			log.Println("Database Query Error ::", err.Error())
			return slots, err
		}
		defer rowsWeekdays.Close()
		for rowsWeekdays.Next() {
			var weekday time.Weekday
			err = rowsWeekdays.Scan(&weekday)
			if err != nil {
				log.Println("Database Scan Error ::", err.Error())
				return slots, err
			}
			weekdays = append(weekdays, weekday)
		}
		slot.Weekdays = weekdays
		slots = append(slots, slot)
	}
	return slots, nil
}

func timeBefore(h1, h2 string) bool {
	h1t, _ := time.Parse("15:04:05", h1)
	h2t, _ := time.Parse("15:04:05", h2)
	return h1t.Before(h2t)
}

// Match slots hours
func match(s1, s2 model.Slot) (bool, model.Slot) {
	slot := model.Slot{}
	found := false
	matchingWeekdays := weekdayMatch(s1.Weekdays, s2.Weekdays)
	if len(matchingWeekdays) != 0 {
		if timeBefore(s1.InitialTime, s2.FinalTime) && timeBefore(s2.InitialTime, s1.FinalTime) {
			if timeBefore(s1.InitialTime, s2.InitialTime) {
				slot.InitialTime = s2.InitialTime
			} else {
				slot.InitialTime = s1.InitialTime
			}
			if timeBefore(s1.FinalTime, s2.FinalTime) {
				slot.FinalTime = s1.FinalTime
			} else {
				slot.FinalTime = s2.FinalTime
			}
			slot.Weekdays = matchingWeekdays
			found = true
		}
	}
	return found, slot
}

// Match slots weekdays
func weekdayMatch(w1, w2 []time.Weekday) []time.Weekday {
	mathingWeekdays := []time.Weekday{}
	for _, w1d := range w1 {
		for _, w2d := range w2 {
			if w1d == w2d {
				mathingWeekdays = append(mathingWeekdays, w1d)
			}
		}
	}
	return mathingWeekdays
}

// Split Slots to 1-hour collection and separated weekdays
func splitSlots(slots []model.Slot) []model.Slot {
	returningSlots := []model.Slot{}
	weekdayTreated := []model.Slot{}
	for _, slot := range slots {
		for {
			if len(slot.Weekdays) > 1 {
				newWeekdays := []time.Weekday{slot.Weekdays[0]}
				slot.Weekdays = append(slot.Weekdays[:0], slot.Weekdays[1:]...)
				newSlot := model.Slot{InitialTime: slot.InitialTime, FinalTime: slot.FinalTime, Weekdays: newWeekdays}
				weekdayTreated = append(weekdayTreated, newSlot)
			} else {
				weekdayTreated = append(weekdayTreated, slot)
				break
			}
		}
	}

	for _, slot := range weekdayTreated {
		for {
			s1i, _ := time.Parse("15:04:05", slot.InitialTime)
			s1f, _ := time.Parse("15:04:05", slot.FinalTime)
			if s1i.Add(time.Duration(1) * time.Hour).Before(s1f) {
				newSlot := model.Slot{}
				newSlot.InitialTime = slot.InitialTime
				newSlot.FinalTime = s1i.Add(time.Duration(1) * time.Hour).Format("15:04:05")
				newSlot.Weekdays = slot.Weekdays
				returningSlots = append(returningSlots, newSlot)
				slot.InitialTime = newSlot.FinalTime
			} else {
				returningSlots = append(returningSlots, slot)
				break
			}
		}
	}
	return returningSlots
}

package model

import "time"

type Candidate struct {
	Id   int    `json:",omitempty"`
	Name string `json:",omitempty"`
}

type Interviewer struct {
	Id   int    `json:",omitempty"`
	Name string `json:",omitempty"`
}

type Slot struct {
	Id          int            `json:",omitempty"`
	PersonId    int            `json:",omitempty"`
	InitialTime string         `json:",omitempty"`
	FinalTime   string         `json:",omitempty"`
	Weekdays    []time.Weekday `json:",omitempty"`
}

type SlotMatchingRequest struct {
	Candidate    Candidate     `json:",omitempty"`
	Interviewers []Interviewer `json:",omitempty"`
}

type SlotMatchingResponse struct {
	Slots []Slot `json:",omitempty"`
}

package structs

import (
	"database/sql/driver"
	"fmt"
)

type InvalidJobState struct {
	State string
}

func (ijs *InvalidJobState) Error() string {
	return fmt.Sprintf("invalid job state %s", ijs.State)
}

type Antling struct {
	Id   int    `json:"id"`
	Jobs []*Job `json:"jobs"`
}

type Antlings struct {
	Antlings []*Antling `json:"antlings"`
}

type JobState int

const (
	JOB_NEW JobState = iota
	JOB_PENDING
	JOB_FINISH
	JOB_ERROR
)

var JOB_STATES = [...]string{
	"NEW",
	"PENDING",
	"FINISH",
	"ERROR",
}

type Job struct {
	Id        int      `json:"id"`
	State     JobState `json:"state"`
	IdAntling int      `json:"-"`
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

func (js *JobState) Scan(value interface{}) error {
	s := string(value.([]uint8)[:])
	for i, state := range JOB_STATES {
		if state == s {
			*js = JobState(i)
			return nil
		}
	}
	return &InvalidJobState{s}
}

func (js JobState) Value() (driver.Value, error) {
	return JOB_STATES[int(js)], nil
}

package structs

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const MAX_RETRIES = 3

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

func (js *JobState) MarshalJSON() ([]byte, error) {
	return json.Marshal(JOB_STATES[int(*js)])
}

func (js *JobState) UnmarshalJSON(data []byte) error {
	state := string(data[1 : len(data)-1])
	for i, s := range JOB_STATES {
		if s == state {
			*js = JobState(i)
			return nil
		}
	}
	return &InvalidJobState{state}
}

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
	Retries   int      `json:"-"`
	IdAntling int      `json:"-"`
	Cmd       []string `json:"command"`
	Env       []string `json:"env"`
	Cwd       string   `json:"cwd"`
	Image     Image    `json:"image"`
}

func (j *Job) SanitizeCwd() string {
	if j.Cwd != "" {
		return j.Cwd
	}
	if j.Image.Cwd != "" {
		return j.Image.Cwd
	}
	return "/"
}

func (j *Job) SanitizeEnv() []string {
	env := j.Image.Env
	if j.Image.Env == nil {
		env = []string{}
	}
	if j.Env != nil {
		env = append(env, j.Env...)
	}
	return env
}

func (j *Job) SanitizeCmd() []string {
	if j.Cmd != nil {
		return j.Cmd
	}
	if j.Image.Cmd != nil {
		return j.Image.Cmd
	}
	return []string{"cat", "/etc/hostname"}
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

type Image struct {
	Id       int      `json:"id"`
	Archive  string   `json:"archive"`
	Hostname string   `json:"hostname"`
	Cmd      []string `json:"command"`
	Env      []string `json:"env"`
	Cwd      string   `json:"cwd"`
}

type Images struct {
	Images []*Image `json:"images"`
}

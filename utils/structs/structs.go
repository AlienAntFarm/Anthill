package structs

type Antling struct {
	Id   int    `json:"id"`
	Jobs []*Job `json:"jobs"`
}

type Antlings struct {
	Antlings []*Antling `json:"antlings"`
}

const (
	JOB_NEW = iota
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
	Id    int `json:"id"`
	State int `json:"state"`
}

type Jobs struct {
	Jobs []*Job `json:"jobs"`
}

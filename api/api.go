package api

import (
	"database/sql"
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"math/rand"
	"time"
)

var Router = mux.NewRouter()
var Scheduler *scheduler

type scheduler struct {
	antlings []int
	queue    map[int]map[int]*structs.Job
	channel  chan *structs.Job
	seed     *rand.Rand
}

func InitScheduler() {
	if Scheduler != nil {
		glog.Fatalf("scheduler already inited, something bad is happening")
	} else {
		Scheduler = newScheduler()
		go Scheduler.start()
	}
}

func newScheduler() *scheduler {
	glog.Infoln("init scheduler")

	id := 0
	antlings := []int{}
	queue := make(map[int]map[int]*structs.Job)
	channel := make(chan *structs.Job, 1)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &scheduler{antlings, queue, channel, seed}

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"
	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			glog.Fatalln(err)
		}
		s.AddAntling(id)
	}

	query = "SELECT id, state, fk_antling "
	query += "FROM anthive.job "
	query += "WHERE fk_antling IS NOT NULL"
	rows, err = db.Conn().Query(query)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		job := &structs.Job{}
		err = rows.Scan(&job.Id, &job.State, &job.IdAntling)
		if err != nil {
			glog.Fatalln(err)
		}
		queue[job.IdAntling][job.Id] = job
		msg := "retrieved job %d from db and assign it to antling %d"
		glog.Infof(msg, job.Id, job.IdAntling)
	}

	return s
}

func (s *scheduler) start() {
	glog.Infoln("starting scheduler")

	queryAssign := "UPDATE anthive.job "
	queryAssign += "SET fk_antling = $1"
	queryAssign += "WHERE anthive.job.id = $2"

	queryUpdate := "UPDATE anthive.job "
	queryUpdate += "SET state = $1 "
	queryUpdate += "WHERE anthive.job.id = $2"

	var row *sql.Row
	for job := range s.channel {
		// if no antling, just pass
		if len(s.antlings) == 0 {
			continue
		}
		if job.IdAntling == 0 {
			// just choose an antling randomly and assign it the job
			job.IdAntling = s.antlings[s.seed.Intn(len(s.antlings))]

			glog.Infof("adding job %d to antling %d", job.Id, job.IdAntling)
			row = db.Conn().QueryRow(queryAssign, job.IdAntling, job.Id)
		} else {
			// if job already assigned we are running updates
			glog.Infof("updating job %d to state %s", job.Id, structs.JOB_STATES[int(job.State)])
			row = db.Conn().QueryRow(queryUpdate, job.State, job.Id)
		}
		err := row.Scan()
		if err != nil && err != sql.ErrNoRows {
			glog.Errorln(err)
			return
		}
		s.queue[job.IdAntling][job.Id] = job

	}
}

func (s *scheduler) ProcessJob(job *structs.Job) {
	s.channel <- job
}

func (s *scheduler) AddAntling(id int) {
	msg := "adding antling %d to cluster, cluster size is now %d"
	s.antlings = append(s.antlings, id)
	glog.Infof(msg, id, len(s.antlings))
	s.queue[id] = make(map[int]*structs.Job)
}

func (s *scheduler) GetJobs(id int) []*structs.Job {
	jobs := []*structs.Job{}
	if jobsMap, ok := s.queue[id]; ok {
		for _, job := range jobsMap {
			jobs = append(jobs, job)
		}
	}
	return jobs
}

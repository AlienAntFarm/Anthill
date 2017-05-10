package api

import (
	"database/sql"
	"github.com/alienantfarm/anthive/db"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"math/rand"
	"time"
)

var Router = mux.NewRouter()
var Scheduler *scheduler

type scheduler struct {
	antlings []int
	queue    map[int][]*Job
	channel  chan int
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
	antlings := []int{}
	queue := make(map[int][]*Job)

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"
	rows, err := db.Conn().Query(query)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var antlingId int
		err = rows.Scan(&antlingId)
		if err != nil {
			glog.Fatalln(err)
		}
		antlings = append(antlings, antlingId)
		queue[antlingId] = []*Job{}
	}

	query = "SELECT id, fk_antling "
	query += "FROM anthive.job "
	query += "WHERE fk_antling IS NOT NULL"
	rows, err = db.Conn().Query(query)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		job := &Job{}
		err = rows.Scan(&job.Id, &id)
		if err != nil {
			glog.Fatalln(err)
		}
		queue[id] = append(queue[id], job)
		msg := "retrieved job %d from db and assign it to antling %d"
		glog.Infof(msg, job.Id, id)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &scheduler{antlings, queue, make(chan int, 1), r}

	return s
}

func (s *scheduler) start() {
	glog.Infoln("starting scheduler")
	query := "UPDATE anthive.job "
	query += "SET fk_antling = $1 "
	query += "WHERE anthive.job.id = $2"

	for job := range s.channel {
		// if no antling, just pass
		if len(s.antlings) == 0 {
			continue
		}
		// just choose an antling randomly and assign it the job
		id := s.antlings[s.seed.Intn(len(s.antlings))]
		glog.Infof("adding job %d to antling %d", job, id)

		row := db.Conn().QueryRow(query, id, job)
		err := row.Scan()
		if err != nil && err != sql.ErrNoRows {
			glog.Errorln(err)
			return
		}

		s.queue[id] = append(s.queue[id], &Job{Id: job})
	}
}

func (s *scheduler) AddJob(id int) {
	s.channel <- id
}

func (s *scheduler) AddAntling(id int) {
	msg := "adding antling %d to cluster, cluster size is now %d"
	s.antlings = append(s.antlings, id)
	glog.Infof(msg, id, len(s.antlings))
	s.queue[id] = []*Job{}
}

func (s *scheduler) GetJobs(id int) []*Job {
	jobs := s.queue[id]
	if jobs != nil {
		return jobs
	} else {
		return []*Job{}
	}
}

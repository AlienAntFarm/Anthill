package api

import (
	"database/sql"
	"github.com/alienantfarm/anthive/db"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/utils/structs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
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

	queryAntlings := "SELECT anthive.antling.id "
	queryAntlings += "FROM anthive.antling"
	rows, err := db.Conn().Query(queryAntlings)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		if rows.Scan(&id); err != nil {
			glog.Fatalln(err)
		}
		s.AddAntling(id)
	}

	queryJobs := "SELECT j.id, j.state, j.cwd, j.command, j.environment, j.fk_antling, "
	queryJobs += "  i.id, i.archive, i.command, i.environment, i.cwd, i.hostname "
	queryJobs += "FROM anthive.job AS j, anthive.image AS i "
	queryJobs += "WHERE j.fk_image = i.id AND j.state < $1"

	rows, err = db.Conn().Query(queryJobs, structs.JOB_FINISH)
	if err != nil {
		glog.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		j := &structs.Job{}
		i := &j.Image

		args := []interface{}{
			&j.Id, &j.State, &j.Cwd, pq.Array(&j.Cmd), pq.Array(&i.Env), &j.IdAntling,
			&i.Id, &i.Archive, pq.Array(&i.Cmd), pq.Array(&i.Env), &i.Cwd, &i.Hostname,
		}
		if err := rows.Scan(args...); err != nil {
			glog.Fatalf("%s", err)
		}
		if _, err := s.schedule(j); err != nil {
			glog.Fatalf("%s", err)
		}
		queue[j.IdAntling][j.Id] = j
		msg := "retrieved job %d from db and assign it to antling %d"
		glog.Infof(msg, j.Id, j.IdAntling)
	}

	return s
}

func (s *scheduler) schedule(job *structs.Job) (sheduled bool, err error) {
	if job.IdAntling != 0 { // we do not want to schedule again
		return
	}
	query := "UPDATE anthive.job "
	query += "SET fk_antling = $1"
	query += "WHERE anthive.job.id = $2"
	update := func(job *structs.Job) error {
		return db.Conn().QueryRow(query, job.IdAntling, job.Id).Scan()
	}

	// just choose an antling randomly and assign it the job
	job.IdAntling = s.antlings[s.seed.Intn(len(s.antlings))]

	if err = update(job); err != nil && err != sql.ErrNoRows {
		job.IdAntling = 0 // reset for further scheduling
		return
	}
	glog.Infof("adding job %d to antling %d", job.Id, job.IdAntling)
	return true, nil
}

func (s *scheduler) start() {
	glog.Infoln("starting scheduler")

	query := "UPDATE anthive.job "
	query += "SET state = $1 "
	query += "WHERE anthive.job.id = $2"

	update := func(job *structs.Job) error {
		return db.Conn().QueryRow(query, job.State, job.Id).Scan()
	}
	for job := range s.channel {
		if len(s.antlings) == 0 { // if no antling, just pass
			continue
		}
		if ok, err := s.schedule(job); ok {
			// job has been scheduled don't do anything
		} else if err != nil {
			glog.Errorf("%s occured when scheduling %d", err, job.Id)
			continue
		} else if err := update(job); err != nil && err != sql.ErrNoRows {
			glog.Errorf("%s occured when updating %d to state %s", err, job.Id, job.State)
			continue
		} else {
			glog.Infof("updating job %d to state %s", job.Id, job.State)
			// job has been updated if state is FINISH or ERROR remove it from scheduler
			if job.State > structs.JOB_PENDING {
				delete(s.queue[job.IdAntling], job.Id)
				continue
			}
		}
		if glog.V(2) {
			glog.Infof(utils.MarshalJSON(job))
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

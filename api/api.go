package api

import (
	"github.com/alienantfarm/anthive/ext/db"
	"github.com/alienantfarm/anthive/ext/minio"
	"github.com/alienantfarm/anthive/utils"
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
	queue    map[int]structs.JobMap
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

	queue := make(map[int]structs.JobMap)
	channel := make(chan *structs.Job, 1)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &scheduler{[]int{}, queue, channel, seed}
	antlings := &db.Antlings{}
	jobs := &db.Jobs{}

	if err := antlings.Get(); err != nil {
		glog.Fatalf("%s", err)
	} else {
		for _, antling := range antlings.Antlings {
			s.AddAntling(antling.Id)
		}
	}

	if err := jobs.Get(structs.JOB_PENDING); err != nil {
		glog.Fatalf("%s", err)
	} else {
		for _, job := range jobs.Jobs {
			queue[job.IdAntling][job.Id] = job
			glog.Infof(
				"retrieved job %d from db and assign it to antling %d",
				job.Id, job.IdAntling,
			)
		}
	}
	return s
}

func (s *scheduler) schedule(job *db.Job) (sheduled bool, err error) {
	if job.IdAntling != 0 { // we do not want to schedule again
		return
	}
	// just choose an antling randomly and assign it the job
	job.IdAntling = s.antlings[s.seed.Intn(len(s.antlings))]
	if minio.Client().MakeBucket("", "us-east-1"); err != nil {
		job.IdAntling = 0
		return
	} else if err = job.UpdateAntling(); err != nil {
		job.IdAntling = 0 // reset for further scheduling
		return
	}
	glog.Infof("adding job %d to antling %d", job.Id, job.IdAntling)
	return true, nil
}

func (s *scheduler) start() {
	glog.Infoln("starting scheduler")

	for job := range s.channel {
		if len(s.antlings) == 0 { // if no antling, just pass
			continue
		}
		if ok, err := s.schedule((*db.Job)(job)); ok {
			// job has been scheduled don't do anything
		} else if err != nil {
			glog.Errorf("%s occured when scheduling %d", err, job.Id)
			continue
		} else if err := ((*db.Job)(job)).UpdateState(); err != nil {
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

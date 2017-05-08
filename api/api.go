package api

import (
	"database/sql"
	"github.com/alienantfarm/anthive/utils"
	"github.com/alienantfarm/anthive/db"
	"github.com/gorilla/mux"
	"math/rand"
	"time"
)

var Router = mux.NewRouter()
var Scheduler = newScheduler()

type scheduler struct {
	antlings []int
	queue    map[int][]int
	channel  chan int
	seed     *rand.Rand
}

func newScheduler() *scheduler {
	antlings := []int{}
	queue := make(map[int][]int)

	query := "SELECT anthive.antling.id "
	query += "FROM anthive.antling"
	rows, err := db.Conn.Query(query)
	if err != nil {
		utils.Error.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var antlingId int
		err = rows.Scan(&antlingId)
		if err != nil {
			utils.Error.Fatalln(err)
		}
		antlings = append(antlings, antlingId)
		queue[antlingId] = []int{}
	}

	query = "SELECT id, fk_antling "
	query += "FROM anthive.job "
	query += "WHERE fk_antling IS NOT NULL"
	rows, err = db.Conn.Query(query)
	if err != nil {
		utils.Error.Fatalln(err)
	}
	defer rows.Close()

	for rows.Next() {
		var antlingId int
		var jobId int
		err = rows.Scan(&jobId, &antlingId)
		if err != nil {
			utils.Error.Fatalln(err)
		}
		queue[antlingId] = append(queue[antlingId], jobId)
		msg := "retrieved job %d from db and assign it to antling %d"
		utils.Info.Printf(msg, jobId, antlingId)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := &scheduler{antlings, queue, make(chan int, 1), r}
	go s.start()

	return s
}

func (s *scheduler) start() {
	query := "UPDATE anthive.job "
	query += "SET fk_antling = $1 "
	query += "WHERE anthive.job.id = $2"

	for jobId := range s.channel {
		// if no antling, just pass
		if len(s.antlings) == 0 {
			continue
		}
		// just choose an antling randomly and assign it the job
		antlingId := s.antlings[s.seed.Intn(len(s.antlings))]
		utils.Info.Printf("adding job %d to antling %d", jobId, antlingId)

		row := db.Conn.QueryRow(query, antlingId, jobId)
		err := row.Scan()
		if err != nil && err != sql.ErrNoRows {
			utils.Error.Printf("%s", err.Error())
			return
		}

		s.queue[antlingId] = append(s.queue[antlingId], jobId)
	}
}

func (s *scheduler) AddJob(jobId int) {
	s.channel <- jobId
}

func (s *scheduler) AddAntling(antlingId int) {
	msg := "adding antling %d to cluster, cluster size is now %d"
	s.antlings = append(s.antlings, antlingId)
	utils.Info.Printf(msg, antlingId, len(s.antlings))
	s.queue[antlingId] = []int{}
}

func (s *scheduler) GetJobs(antlingId int) []int {
	jobs := s.queue[antlingId]
	if jobs != nil {
		return jobs
	} else {
		return []int{}
	}
}

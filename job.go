package kuntul

import (
	"github.com/robfig/cron/v3"
)

type Job struct {
	adapter Adapter
	cron    *cron.Cron
}

func NewJob(adapter Adapter) *Job {
	return &Job{
		adapter: adapter,
		cron:    cron.New(cron.WithSeconds()),
	}
}

func (j *Job) Add(task *Task) error {
	job := cron.FuncJob(func() {
		err := j.adapter.Lock(task)
		if err == nil {
			task.Cmd(j)
		}
	})

	j.cron.AddJob(task.Schedule, job)

	return nil
}

func (j *Job) Start() {
	go j.cron.Start()
}

func (j *Job) Done() error {
	return j.adapter.Unlock()
}

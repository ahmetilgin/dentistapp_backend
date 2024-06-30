package localcache

import (
	"backend/job"
	"context"
	"models"
	"sync"
)

type JobLocalStorage struct {
	jobs map[string]*models.Job
	mutex     *sync.Mutex
}

func NewJobLocalStorage() *JobLocalStorage {
	return &JobLocalStorage{
		jobs: make(map[string]*models.Job),
		mutex:     new(sync.Mutex),
	}
}

func (s *JobLocalStorage) CreateJob(ctx context.Context, user *models.User, bm *models.Job) error {
	bm.UserID = user.ID

	s.mutex.Lock()
	s.jobs[bm.ID] = bm
	s.mutex.Unlock()

	return nil
}

func (s *JobLocalStorage) GetJobs(ctx context.Context) ([]*models.Job, error) {
	jobs := make([]*models.Job, 0)

	s.mutex.Lock()
	for _, bm := range s.jobs {
		jobs = append(jobs, bm)
	}
	s.mutex.Unlock()

	return jobs, nil
}

func (s *JobLocalStorage) DeleteJob(ctx context.Context, user *models.User, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	bm, ex := s.jobs[id]
	if ex && bm.UserID == user.ID {
		delete(s.jobs, id)
		return nil
	}

	return job.ErrJobNotFound
}

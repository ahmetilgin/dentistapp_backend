package mock

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type JobStorageMock struct {
	mock.Mock
}

func (s *JobStorageMock) CreateJob(ctx context.Context, user *models.User, bm *models.Job) error {
	args := s.Called(user, bm)

	return args.Error(0)
}

func (s *JobStorageMock) GetJobs(ctx context.Context) ([]*models.Job, error) {
	args := s.Called()

	return args.Get(0).([]*models.Job), args.Error(1)
}

func (s *JobStorageMock) DeleteJob(ctx context.Context, user *models.User, id string) error {
	args := s.Called(user, id)

	return args.Error(0)
}

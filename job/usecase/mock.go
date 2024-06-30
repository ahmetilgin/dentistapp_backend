package usecase

import (
	"backend/models"
	"context"

	"github.com/stretchr/testify/mock"
)

type JobUseCaseMock struct {
	mock.Mock
}

func (m JobUseCaseMock) CreateJob(ctx context.Context, user *models.User, job * models.Job) error {
	args := m.Called(user, job)

	return args.Error(0)
}

func (m JobUseCaseMock) GetJob(ctx context.Context, user *models.User) ([]*models.Job, error) {
	args := m.Called(user)

	return args.Get(0).([]*models.Job), args.Error(1)
}

func (m JobUseCaseMock) DeleteJob(ctx context.Context, user *models.User, id string) error {
	args := m.Called(user, id)

	return args.Error(0)
}

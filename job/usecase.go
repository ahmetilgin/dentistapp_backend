package job

import (
	"backend/models"
	"context"
)

type UseCase interface {
	CreateJob(ctx context.Context, user *models.User, job* models.Job) error
	GetJobs(ctx context.Context) ([]*models.Job, error)
	DeleteJob(ctx context.Context, user *models.User, id string) error
}

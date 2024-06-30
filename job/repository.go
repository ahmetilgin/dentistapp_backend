package job

import (
	"backend/models"
	"context"
)

type Repository interface {
	CreateJob(ctx context.Context, user *models.User, bm *models.Job) error
	GetJobs(ctx context.Context) ([]*models.Job, error)
	DeleteJob(ctx context.Context, user *models.User, id string) error
}

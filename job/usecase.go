package job

import (
	"backend/models"
	"context"
)

type UseCase interface {
	CreateJob(ctx context.Context, user *models.BusinessUser, job* models.Job) error
	GetJobs(ctx context.Context) ([]*models.Job, error)
	Search(ctx context.Context, location, keyword string) ([]*models.Job, error)
	SearchProfession(ctx context.Context, keyword string) ([]*models.Profession, error)
	DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error
}

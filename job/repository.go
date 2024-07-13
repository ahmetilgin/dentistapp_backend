package job

import (
	"backend/models"
	"context"
)

type Repository interface {
	CreateJob(ctx context.Context, user *models.BusinessUser, bm *models.Job) error
	GetJobs(ctx context.Context) ([]*models.Job, error)
	DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error
	Search(ctx context.Context, location, keyword string) ([]*models.Job, error)
	SearchProfession(ctx context.Context,  keyword string) ([]*models.Profession, error)
	GetPopulerJobs(ctx context.Context) ([]*models.Profession, error)
}

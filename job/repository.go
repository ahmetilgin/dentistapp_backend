package job

import (
	"backend/models"
	"context"
)

type Repository interface {
	CreateJob(ctx context.Context, user *models.BusinessUser, bm *models.Job) error
	GetJobs(ctx context.Context) ([]*models.Job, error)
	DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error
	Search(ctx context.Context, location, keyword, region string) ([]*models.Job, error)
	SearchProfession(ctx context.Context, keyword, region string) ([]*models.Profession, error)
	GetPopulerJobs(ctx context.Context, code string) ([]*models.Profession, error)
}

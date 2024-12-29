package job

import (
	jobmongo "backend/job/repository/mongo"
	"backend/models"
	"context"
)

type UseCase interface {
	CreateJob(ctx context.Context, user *models.BusinessUser, job *models.Job) error
	Search(ctx context.Context, position, region string, page, limit int) ([]*jobmongo.JobDetails, error)
	SearchProfession(ctx context.Context, keyword, region string) ([]string, error)
	ApplyJob(ctx context.Context, user *models.NormalUser, jobId string) error
	DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error
	GetPopulerJobs(ctx context.Context, region string) ([]string, error)
	GetJobs(ctx context.Context, user *models.BusinessUser) ([]*models.Job, error)
	Update(ctx context.Context, user *models.BusinessUser, job *models.Job) error
}

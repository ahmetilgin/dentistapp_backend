package usecase

import (
	"backend/job"
	"backend/models"
	"context"
)

type JobUseCase struct {
	jobRepo job.Repository
}

func NewJobUseCase(jobRepo job.Repository) *JobUseCase {
	return &JobUseCase{
		jobRepo: jobRepo,
	}
}

func (b JobUseCase) CreateJob(ctx context.Context, user *models.BusinessUser, job *models.Job ) error {
	return b.jobRepo.CreateJob(ctx, user, job)
}

func (b JobUseCase) GetJobs(ctx context.Context) ([]*models.Job, error) {
	return b.jobRepo.GetJobs(ctx)
}

func (b JobUseCase) DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error {
	return b.jobRepo.DeleteJob(ctx, user, id)
}

func (b JobUseCase) Search(ctx context.Context, location, keyword string) ([]*models.Job, error) {
	return b.jobRepo.Search(ctx, location, keyword)
}

func (b JobUseCase) SearchProfession(ctx context.Context, keyword string) ([]*models.Profession, error){
	return b.jobRepo.SearchProfession(ctx, keyword)
}




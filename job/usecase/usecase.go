package usecase

import (
	"backend/job"
	jobmongo "backend/job/repository/mongo"
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

func (b JobUseCase) CreateJob(ctx context.Context, user *models.BusinessUser, job *models.Job) error {
	return b.jobRepo.CreateJob(ctx, user, job)
}

func (b JobUseCase) DeleteJob(ctx context.Context, user *models.BusinessUser, id string) error {
	return b.jobRepo.DeleteJob(ctx, user, id)
}

func (b JobUseCase) Search(ctx context.Context, location, keyword, region string) ([]*jobmongo.JobDetails, error) {
	return b.jobRepo.Search(ctx, location, keyword, region)
}

func (b JobUseCase) SearchProfession(ctx context.Context, keyword, region string) ([]string, error) {
	professions, _ := b.jobRepo.SearchProfession(ctx, keyword, region)
	var queryResult []string

	for _, profession := range professions {
		queryResult = append(queryResult, profession.Name)
	}

	return queryResult, nil
}

func (b JobUseCase) GetPopulerJobs(ctx context.Context, code string) ([]string, error) {
	professions, _ := b.jobRepo.GetPopulerJobs(ctx, code)
	var queryResult []string

	for _, profession := range professions {
		queryResult = append(queryResult, profession.Name)
	}

	return queryResult, nil
}

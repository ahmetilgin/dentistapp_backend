package usecase

import (
	authmongo "backend/auth/repository/mongo"
	"backend/job"
	jobmongo "backend/job/repository/mongo"
	"backend/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (b JobUseCase) Search(ctx context.Context, location, keyword string, page, limit int) ([]*jobmongo.JobDetails, error) {
	return b.jobRepo.Search(ctx, location, keyword, page, limit)
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

func (b JobUseCase) ApplyJob(ctx context.Context, user *models.NormalUser, jobId string) error {
	return b.jobRepo.ApplyJob(ctx, user, jobId)
}

func (b JobUseCase) GetJobs(ctx context.Context, user *models.BusinessUser) ([]*models.Job, error) {
	return b.jobRepo.GetJobs(ctx, user)
}

func (b JobUseCase) Update(ctx context.Context, user *models.BusinessUser, job *models.Job) error {
	return b.jobRepo.Update(ctx, user, job)
}

func (b JobUseCase) GetCandidateDetails(ctx context.Context, user *models.BusinessUser, candidateId string) (*models.NormalUser, error) {
	// Verify that the user owns the job that this candidate applied to
	_, err := b.jobRepo.GetJobs(ctx, user)
	if err != nil {
		return nil, err
	}

	// Get the user repository from the job repository to fetch user details
	userRepo := b.jobRepo.GetUserRepository()

	// Fetch candidate details using the concrete implementation
	userRepoMongo, ok := userRepo.(*authmongo.UserRepository)
	if !ok {
		return nil, fmt.Errorf("failed to convert user repository")
	}

	candidateIdPrimitive, err := primitive.ObjectIDFromHex(candidateId)
	if err != nil {
		return nil, err
	}

	candidate, err := userRepoMongo.GetNormalUserById(ctx, candidateIdPrimitive)
	if err != nil {
		return nil, err
	}

	return candidate, nil
}

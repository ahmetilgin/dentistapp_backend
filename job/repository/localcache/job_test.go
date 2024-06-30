package localcache

import (
	"backend/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	// id := "id"
	// user := &models.User{ID: id}

	// // s := NewJobLocalStorage()

	// // for i := 0; i < 10; i++ {
	// // 	bm := &models.Job{
	// // 		ID:     fmt.Sprintf("id%d", i),
	// // 		UserID: user.ID,
	// // 	}

	// // 	err := s.CreateJob(context.Background(), user, bm)
	// // 	assert.NoError(t, err)
	// // }

	// returnedJobs, err := s.GetJobs(context.Background())
	// assert.NoError(t, err)

	// assert.Equal(t, 10, len(returnedJobs))
}

func TestDeleteJob(t *testing.T) {
	id1 := "id1"
	id2 := "id2"

	user1 := &models.User{ID: id1}
	user2 := &models.User{ID: id2}

	bmID := "bmID"
	bm := &models.Job{ID: bmID, UserID: user1.ID}

	s := NewJobLocalStorage()

	err := s.CreateJob(context.Background(), user1, bm)
	assert.NoError(t, err)

	err = s.DeleteJob(context.Background(), user1, bmID)
	assert.NoError(t, err)

	err = s.CreateJob(context.Background(), user1, bm)
	assert.NoError(t, err)

	err = s.DeleteJob(context.Background(), user2, bmID)
	assert.Error(t, err)
	assert.Equal(t, err, Job.ErrJobNotFound)
}

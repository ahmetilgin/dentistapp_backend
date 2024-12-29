package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Professions struct {
	Code        string       `bson:"code"`
	Professions []Profession `bson:"professions"`
}

type Profession struct {
	Name          string `bson:"name"`
	SearchCounter int    `bson:"count"`
}

type Job struct {
	ID                  primitive.ObjectID   `bson:"_id,omitempty"`
	UserID              primitive.ObjectID   `bson:"user_id,omitempty"`
	JobTitle            string               `bson:"job_title,omitempty"`
	Description         string               `bson:"description,omitempty"`
	Location            string               `bson:"location,omitempty"`
	SalaryRange         string               `bson:"salary_range,omitempty"`
	EmploymentType      string               `bson:"employment_type,omitempty"` // full-time, part-time, contract, etc.
	DatePosted          string               `bson:"date_posted,omitempty"`
	ApplicationDeadline string               `bson:"application_deadline,omitempty"`
	Candidates          []primitive.ObjectID `bson:"candidates,omitempty"`
}

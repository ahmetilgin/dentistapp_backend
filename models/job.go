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
	UserID              primitive.ObjectID `bson:"user_id,omitempty"`
	JobTitle            string             `bson:"job_title,omitempty"`
	Description         string             `bson:"description,omitempty"`
	Requirements        string             `bson:"requirements,omitempty"`
	Location            string             `bson:"location,omitempty"`
	SalaryRange         string             `bson:"salary_range,omitempty"`
	EmploymentType      string             `bson:"employment_type,omitempty"` // full-time, part-time, contract, etc.
	DatePosted          string             `bson:"date_posted,omitempty"`
	ApplicationDeadline string             `bson:"application_deadline,omitempty"`
}

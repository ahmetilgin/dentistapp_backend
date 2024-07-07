package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profession struct {
    Name         	string 	`bson:"name"`
	SearchCounter 	int   	`bson:"count"`
}

type Job struct {
	UserID 				primitive.ObjectID `bson:"user_id"`
	JobTitle            string `bson:"job_title,omitempty"`
	Description         string `bson:"description,omitempty"`
	Requirements        string `bson:"requirements,omitempty"`  
	Location            string `bson:"location,omitempty"` 
	SalaryRange         string `bson:"salary_range,omitempty"` 
	EmploymentType      string `bson:"employment_type,omitempty"`  // full-time, part-time, contract, etc.
	DatePosted          time.Time `bson:"date_posted,omitempty"`  
	ApplicationDeadline time.Time `bson:"application_deadline,omitempty"`  
}
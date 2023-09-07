package models

import (
	"time"
)

type TestCategory struct {
	ID       int    `json:"id,omitempty" bson:"id,omitempty"`
	Category string `json:"category" bson:"category"`
}

type Dashboard struct {
	ID                int           `json:"id,omitempty" bson:"id,omitempty"`
	Category          string        `json:"category" bson:"category"`
	Description       string        `json:"description" bson:"description"`
	NumberOfQuestions int           `json:"numberofQuestions" bson:"numberofQuestions"`
	Duration          time.Duration `json:"duration" bson:"duration"`
	// Category          TestCategory  `json:"testcategory" bson:"testcategory"`
	// Status            string        `json:"status" bson:"status"`
	// LastModified      time.Time     `json:"lastmodified" bson:"lastmodified"`
}

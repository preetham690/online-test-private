package models

type UserResult struct {
	Email  string `json:"email" bson:"email"`
	TestID int    `json:"testID" bson:"testID"`
	Score  int    `json:"score" bson:"score"`
}

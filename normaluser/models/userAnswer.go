package models

type UserAnswer struct {
	Email           string           `json:"email" bson:"email"`
	TestID          int              `json:"testID" bson:"testID"`
	QuestionAnswers []QuestionAnswer `json:"questionAnswers" bson:"questionAnswers"`
}

type QuestionAnswer struct {
	QuestionID int    `json:"questionID" bson:"questionID"`
	Answer     string `json:"answer" bson:"answer"`
}

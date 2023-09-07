package models

// Define the TestDetails structure
type TestDetails struct {
	TestID    int        `json:"testID" bson:"testID"`
	Category  string     `json:"category" bson:"category"`
	Questions []Question `json:"questions" bson:"questions"`
	Weightage int        `json:"weightage" bson:"weightage"`
}

// Define the Question structure
type Question struct {
	QuestionID int      `json:"questionID" bson:"questionID"`
	Text       string   `json:"text" bson:"text"`
	Options    []string `json:"options" bson:"options"`
	Answer     string   `json:"answer" bson:"answer"`
}

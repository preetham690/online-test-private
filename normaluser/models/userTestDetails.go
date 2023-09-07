package models

type UserTestDetails struct {
	TestID int `json:"testID" bson:"testID"`
	//Questions []Question `json:"question" bson:"question"`
	Category string `json:"category" bson:"category"`
}

// Define the Question structure
type Question struct {
	QuestionID int      `json:"questionID" bson:"questionID"`
	Text       string   `json:"text" bson:"text"`
	Options    []string `json:"options" bson:"options"`
	Answer     string   `json:"answer" bson:"answer"`
}

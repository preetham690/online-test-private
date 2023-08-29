package models

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	User_ID  primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Mobile   string             `json:"mobile" bson:"mobile"`
	Password string             `json:"password" bson:"password"`
	IsAdmin  bool               `json:"isAdmin" bson:"isAdmin"`
}

type LoginForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("invalid email address field")
	}
	if u.Password == "" {
		return errors.New("invalid password field")
	}
	if u.Mobile == "" {
		return errors.New("invalid mobile number")
	}
	return nil
}

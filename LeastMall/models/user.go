package models

type User struct {
	Id       int
	Phone    string
	Password string
	Email    string
	Status   int
}

func (u User) TableName() string {
	return "user"
}

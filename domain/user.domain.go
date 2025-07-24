package domain

type User struct {
	ID    string
	Name  string
	Email string
	Img   *string
	Level int8
}

type UserRequest struct {
	Id   string
	Name string
	Img  string
}

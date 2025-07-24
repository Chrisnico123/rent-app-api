package web

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Img   string `json:"img"`
	Level int8   `json:"level"`
}

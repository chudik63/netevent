package models

type User struct {
	Id        int64  `json: "id" db: "id"`
	Name      string `json: "name" db: "name"`
	Email     string `json: "email" db: "email"`
	Password  string `json: "password" db: "password"`
	Role      string `json: "role" db: "role"`
	Interests string `json: "interest" db: "interest"`
}

type Token struct {
	AccessTkn  string `json: "accesstkn" db: "accesstkn"`
	AccessTtl  int64  `json: "accessttl" db: "accessttl"`
	RefreshTkn string `json: "refreshtkn" db: "refreshtkn"`
	RefreshTtl int64  `json: "refreshttl" db: "refreshttl"`
}

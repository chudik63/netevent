package models

/*
	message User {
	  int64 id = 1;
	  string name = 2;
	  string email = 3;
	  string password = 4; // Should be hashed before storing
	  repeated string interests = 5;
	}

	message Token {
	  string access_token = 1;
	  int64 access_token_ttl = 2; // Time in seconds
	  string refresh_token = 3;
	  int64 refresh_token_ttl = 4; // Time in seconds
	}

CREATE TABLE IF NOT EXISTS "tuser" (

	id INT PRIMARY KEY NOT NULL,
	name VARCHAR(30),
	password VARCHAR(30),
	email   VARCHAR(30),
	role    int,
	interest TEXT,
	accesstkn TEXT,
	accessttl INT,
	refreshtkn TEXT,
	refreshttl INT

);
*/

type User struct {
	Id        int64  `json: "id" db: "id"`
	Name      string `json: "name" db: "name"`
	Email     string `json: "email" db: "email"`
	Password  string `json: "password" db: "password"`
	Interests string `json: "interest" db: "interest"`
}

type Token struct {
	AccessTkn  string `json: "accesstkn" db: "accesstkn"`
	AccessTtl  int64  `json: "accessttl" db: "accessttl"`
	RefreshTkn string `json: "refreshtkn" db: "refreshtkn"`
	RefreshTtl int64  `json: "refreshttl" db: "refreshttl"`
}

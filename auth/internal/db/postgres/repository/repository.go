package repository

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres"
	"gitlab.crja72.ru/gospec/go9/netevent/auth/internal/db/postgres/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db.Db}
}

/*
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
*/
func (u *UserRepository) NewUser(users *models.User) error {
	_, err := sq.Insert("tuser").Columns("id", "name", "password", "email", "interests").
		Values(users.Id, users.Name, users.Password, users.Email, users.Interests).
		RunWith(u.db).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) AuthUser(name string, password string, tkn *models.Token) error {
	us := models.User{}
	err := sq.Select("id", "name", "password").
		From("tuser").
		Where(sq.And{sq.Eq{"name": name}, sq.Eq{"password": password}}).
		RunWith(u.db).QueryRow().
		Scan(&us.Id, &us.Name, &us.Password)

	if err != nil {
		return err
	}
	log.Println(us)

	res, err := sq.Update("tuser").
		Set("accesstkn", tkn.AccessTkn).Set("accessttl", tkn.AccessTtl).
		Set("refreshtkn", tkn.RefreshTkn).Set("refreshttl", tkn.RefreshTtl).
		RunWith(u.db).Exec()
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func (u *UserRepository) UpdateTocken(users *models.User) error {
	return nil
}

func (u *UserRepository) GetId(name string) (int, error) {
	us := models.User{}
	err := sq.Select("id").
		From("tuser").Where(sq.Eq{"name": name}).
		RunWith(u.db).QueryRow().Scan(&us.Id, &us.Password)

	if err != nil {
		return 0, err
	}
	log.Println(us)
	return int(us.Id), nil
}

func (u *UserRepository) GetInterests(id int) (string, error) {
	interests := ""
	err := sq.Select("interest").
		From("tuser").Where(sq.Eq{"id": id}).
		RunWith(u.db).QueryRow().Scan(&interests)
	if err != nil {
		return "", err
	}
	fmt.Println(interests)
	return interests, nil
}

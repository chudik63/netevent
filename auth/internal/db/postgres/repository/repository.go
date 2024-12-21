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

func (u *UserRepository) NewUser(users *models.User) error {
	_, err := sq.Insert("tuser").Columns("id", "name", "password", "email", "role", "interest").
		Values(users.Id, users.Name, users.Password, users.Email, users.Role, users.Interests).PlaceholderFormat(sq.Dollar).
		RunWith(u.db).Exec()
	if err != nil {
		log.Println(err)
		return fmt.Errorf("err in database: %w", err)
	}
	return nil
}

func (u *UserRepository) AuthUser(name string, password string, tkn *models.Token) error {
	us := models.User{}
	err := sq.Select("id", "name", "password").
		From("tuser").
		Where(sq.Eq{"name": name}, sq.Eq{"password": password}).PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRow().
		Scan(&us.Id, &us.Name, &us.Password)

	if err != nil {
		return fmt.Errorf("err in database: %w", err)
	}
	log.Println(us)

	res, err := sq.Update("tuser").
		Set("accesstkn", tkn.AccessTkn).Set("accessttl", tkn.AccessTtl).
		Set("refreshtkn", tkn.RefreshTkn).Set("refreshttl", tkn.RefreshTtl).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).Exec()
	if err != nil {
		return fmt.Errorf("err in database: %w", err)
	}
	fmt.Println(res)
	return nil
}

func (u *UserRepository) UpdateToken(users *models.User) error {
	return nil
}

func (u *UserRepository) GetRole(name string) (string, error) {
	var role, nm string

	err := sq.Select("role", "name").
		From("tuser").Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRow().Scan(&role, &nm)

	if err != nil {
		return "", fmt.Errorf("err in database: %w", err)
	}
	return role, nil
}

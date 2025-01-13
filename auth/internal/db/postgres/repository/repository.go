package repository

import (
	"fmt"
	"log"

	"github.com/chudik63/netevent/auth/internal/db/postgres"
	"github.com/chudik63/netevent/auth/internal/db/postgres/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db.Db}
}

func (u *UserRepository) NewUser(users *models.User) (int64, error) {
	var id int64

	err := sq.Insert("tuser").
		Columns("id", "name", "password", "email", "role", "interest").
		Values(users.Name, users.Password, users.Email, users.Role, users.Interests).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRow().
		Scan(&id)

	if err != nil {
		log.Println(err)
		return 0, fmt.Errorf("err in database: %w", err)
	}

	return id, nil
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

	_, err = sq.Update("tuser").
		Set("accesstkn", tkn.AccessTkn).Set("accessttl", tkn.AccessTtl).
		Set("refreshtkn", tkn.RefreshTkn).Set("refreshttl", tkn.RefreshTtl).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).Exec()
	if err != nil {
		return fmt.Errorf("err in database: %w", err)
	}

	return nil
}

func (u *UserRepository) UpdateToken(users *models.User) error {
	return nil
}

func (u *UserRepository) GetRole(name string) (string, error) {
	var role string

	err := sq.Select("role").
		From("tuser").Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRow().Scan(&role)

	if err != nil {
		return "", fmt.Errorf("err in database: %w", err)
	}
	return role, nil
}

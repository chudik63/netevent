package repository

import (
	"database/sql"
	"errors"

	"github.com/chudik63/netevent/auth_service/internal/db/postgres"
	"github.com/chudik63/netevent/auth_service/internal/models"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{db.Db}
}

func (u *UserRepository) NewUser(user *models.User) (int64, error) {
	var id int64

	err := sq.Select("id").
		From("tuser").
		Where(sq.Eq{"name": user}).
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRow().
		Scan(&id)

	if id != 0 {
		return 0, models.ErrUserAlreadyExists
	}

	err = sq.Insert("tuser").
		Columns("name", "password", "email", "role", "interests").
		Values(user.Name, user.Password, user.Email, user.Role, user.Interests).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		RunWith(u.db).
		QueryRow().
		Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u *UserRepository) AuthUser(name string, password string) (*models.User, error) {
	var user models.User

	err := sq.Select("id", "name", "password", "role", "email").
		From("tuser").
		Where(sq.Eq{"name": name}, sq.Eq{"password": password}).PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRow().
		Scan(&user.Id, &user.Name, &user.Password, &user.Role, &user.Email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

package repository

import (
	"fmt"
	"log"

	"github.com/chudik63/netevent/auth_service/internal/db/postgres"
	"github.com/chudik63/netevent/auth_service/internal/db/postgres/models"

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
		Columns("name", "password", "email", "role", "interests").
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

func (u *UserRepository) AuthUser(name string, password string) (*models.User, error) {
	var user models.User

	err := sq.Select("id", "name", "password", "role", "email").
		From("tuser").
		Where(sq.Eq{"name": name}, sq.Eq{"password": password}).PlaceholderFormat(sq.Dollar).
		RunWith(u.db).QueryRow().
		Scan(&user.Id, &user.Name, &user.Password, &user.Role, &user.Email)

	if err != nil {
		return nil, fmt.Errorf("err in database: %w", err)
	}

	// _, err = sq.Update("tuser").
	// 	Set("accesstkn", tkn.AccessTkn).Set("accessttl", tkn.AccessTtl).
	// 	Set("refreshtkn", tkn.RefreshTkn).Set("refreshttl", tkn.RefreshTtl).
	// 	PlaceholderFormat(sq.Dollar).
	// 	RunWith(u.db).Exec()
	// if err != nil {
	// 	return fmt.Errorf("err in database: %w", err)
	// }

	return &user, nil
}
